// Test sync/copy/move

package fs_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/ncw/rclone/fs"
	"github.com/ncw/rclone/fstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/unicode/norm"
)

// Check dry run is working
func TestCopyWithDryRun(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteFile("sub dir/hello world", "hello world", t1)
	r.Mkdir(r.Fremote)

	fs.Config.DryRun = true
	err := fs.CopyDir(r.Fremote, r.Flocal)
	fs.Config.DryRun = false
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote)
}

// Now without dry run
func TestCopy(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteFile("sub dir/hello world", "hello world", t1)
	r.Mkdir(r.Fremote)

	err := fs.CopyDir(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file1)
}

// Now with --no-traverse
func TestCopyNoTraverse(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()

	fs.Config.NoTraverse = true
	defer func() { fs.Config.NoTraverse = false }()

	file1 := r.WriteFile("sub dir/hello world", "hello world", t1)

	err := fs.CopyDir(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file1)
}

// Now with --no-traverse
func TestSyncNoTraverse(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()

	fs.Config.NoTraverse = true
	defer func() { fs.Config.NoTraverse = false }()

	file1 := r.WriteFile("sub dir/hello world", "hello world", t1)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file1)
}

// Test copy with depth
func TestCopyWithDepth(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteFile("sub dir/hello world", "hello world", t1)
	file2 := r.WriteFile("hello world2", "hello world2", t2)

	// Check the MaxDepth too
	fs.Config.MaxDepth = 1
	defer func() { fs.Config.MaxDepth = -1 }()

	err := fs.CopyDir(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file1, file2)
	fstest.CheckItems(t, r.Fremote, file2)
}

// Test a server side copy if possible, or the backup path if not
func TestServerSideCopy(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteObject("sub dir/hello world", "hello world", t1)
	fstest.CheckItems(t, r.Fremote, file1)

	FremoteCopy, _, finaliseCopy, err := fstest.RandomRemote(*fstest.RemoteName, *fstest.SubDir)
	require.NoError(t, err)
	defer finaliseCopy()
	t.Logf("Server side copy (if possible) %v -> %v", r.Fremote, FremoteCopy)

	err = fs.CopyDir(FremoteCopy, r.Fremote)
	require.NoError(t, err)

	fstest.CheckItems(t, FremoteCopy, file1)
}

// Check that if the local file doesn't exist when we copy it up,
// nothing happens to the remote file
func TestCopyAfterDelete(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteObject("sub dir/hello world", "hello world", t1)
	fstest.CheckItems(t, r.Flocal)
	fstest.CheckItems(t, r.Fremote, file1)

	err := fs.Mkdir(r.Flocal, "")
	require.NoError(t, err)

	err = fs.CopyDir(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal)
	fstest.CheckItems(t, r.Fremote, file1)
}

// Check the copy downloading a file
func TestCopyRedownload(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteObject("sub dir/hello world", "hello world", t1)
	fstest.CheckItems(t, r.Fremote, file1)

	err := fs.CopyDir(r.Flocal, r.Fremote)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file1)
}

// Create a file and sync it. Change the last modified date and resync.
// If we're only doing sync by size and checksum, we expect nothing to
// to be transferred on the second sync.
func TestSyncBasedOnCheckSum(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	fs.Config.CheckSum = true
	defer func() { fs.Config.CheckSum = false }()

	file1 := r.WriteFile("check sum", "", t1)
	fstest.CheckItems(t, r.Flocal, file1)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	// We should have transferred exactly one file.
	assert.Equal(t, int64(1), fs.Stats.GetTransfers())
	fstest.CheckItems(t, r.Fremote, file1)

	// Change last modified date only
	file2 := r.WriteFile("check sum", "", t2)
	fstest.CheckItems(t, r.Flocal, file2)

	fs.Stats.ResetCounters()
	err = fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	// We should have transferred no files
	assert.Equal(t, int64(0), fs.Stats.GetTransfers())
	fstest.CheckItems(t, r.Flocal, file2)
	fstest.CheckItems(t, r.Fremote, file1)
}

// Create a file and sync it. Change the last modified date and the
// file contents but not the size.  If we're only doing sync by size
// only, we expect nothing to to be transferred on the second sync.
func TestSyncSizeOnly(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	fs.Config.SizeOnly = true
	defer func() { fs.Config.SizeOnly = false }()

	file1 := r.WriteFile("sizeonly", "potato", t1)
	fstest.CheckItems(t, r.Flocal, file1)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	// We should have transferred exactly one file.
	assert.Equal(t, int64(1), fs.Stats.GetTransfers())
	fstest.CheckItems(t, r.Fremote, file1)

	// Update mtime, md5sum but not length of file
	file2 := r.WriteFile("sizeonly", "POTATO", t2)
	fstest.CheckItems(t, r.Flocal, file2)

	fs.Stats.ResetCounters()
	err = fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	// We should have transferred no files
	assert.Equal(t, int64(0), fs.Stats.GetTransfers())
	fstest.CheckItems(t, r.Flocal, file2)
	fstest.CheckItems(t, r.Fremote, file1)
}

// Create a file and sync it. Keep the last modified date but change
// the size.  With --ignore-size we expect nothing to to be
// transferred on the second sync.
func TestSyncIgnoreSize(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	fs.Config.IgnoreSize = true
	defer func() { fs.Config.IgnoreSize = false }()

	file1 := r.WriteFile("ignore-size", "contents", t1)
	fstest.CheckItems(t, r.Flocal, file1)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	// We should have transferred exactly one file.
	assert.Equal(t, int64(1), fs.Stats.GetTransfers())
	fstest.CheckItems(t, r.Fremote, file1)

	// Update size but not date of file
	file2 := r.WriteFile("ignore-size", "longer contents but same date", t1)
	fstest.CheckItems(t, r.Flocal, file2)

	fs.Stats.ResetCounters()
	err = fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	// We should have transferred no files
	assert.Equal(t, int64(0), fs.Stats.GetTransfers())
	fstest.CheckItems(t, r.Flocal, file2)
	fstest.CheckItems(t, r.Fremote, file1)
}

func TestSyncIgnoreTimes(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteBoth("existing", "potato", t1)
	fstest.CheckItems(t, r.Fremote, file1)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	// We should have transferred exactly 0 files because the
	// files were identical.
	assert.Equal(t, int64(0), fs.Stats.GetTransfers())

	fs.Config.IgnoreTimes = true
	defer func() { fs.Config.IgnoreTimes = false }()

	fs.Stats.ResetCounters()
	err = fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	// We should have transferred exactly one file even though the
	// files were identical.
	assert.Equal(t, int64(1), fs.Stats.GetTransfers())

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file1)
}

func TestSyncIgnoreExisting(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteFile("existing", "potato", t1)

	fs.Config.IgnoreExisting = true
	defer func() { fs.Config.IgnoreExisting = false }()

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file1)

	// Change everything
	r.WriteFile("existing", "newpotatoes", t2)
	fs.Stats.ResetCounters()
	err = fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)
	// Items should not change
	fstest.CheckItems(t, r.Fremote, file1)
}

func TestSyncAfterChangingModtimeOnly(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteFile("empty space", "", t2)
	file2 := r.WriteObject("empty space", "", t1)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file2)

	fs.Config.DryRun = true
	defer func() { fs.Config.DryRun = false }()

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file2)

	fs.Config.DryRun = false

	fs.Stats.ResetCounters()
	err = fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file1)
}

func TestSyncAfterChangingModtimeOnlyWithNoUpdateModTime(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()

	if r.Fremote.Hashes().Count() == 0 {
		t.Logf("Can't check this if no hashes supported")
		return
	}

	fs.Config.NoUpdateModTime = true
	defer func() {
		fs.Config.NoUpdateModTime = false
	}()

	file1 := r.WriteFile("empty space", "", t2)
	file2 := r.WriteObject("empty space", "", t1)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file2)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file2)
}

func TestSyncDoesntUpdateModtime(t *testing.T) {
	if fs.Config.ModifyWindow == fs.ModTimeNotSupported {
		t.Skip("Can't run this test on fs which doesn't support mod time")
	}
	r := fstest.NewRun(t)
	defer r.Finalise()

	file1 := r.WriteFile("foo", "foo", t2)
	file2 := r.WriteObject("foo", "bar", t1)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file2)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file1)

	// We should have transferred exactly one file, not set the mod time
	assert.Equal(t, int64(1), fs.Stats.GetTransfers())
}

func TestSyncAfterAddingAFile(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteBoth("empty space", "", t2)
	file2 := r.WriteFile("potato", "------------------------------------------------------------", t3)

	fstest.CheckItems(t, r.Flocal, file1, file2)
	fstest.CheckItems(t, r.Fremote, file1)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Flocal, file1, file2)
	fstest.CheckItems(t, r.Fremote, file1, file2)
}

func TestSyncAfterChangingFilesSizeOnly(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteObject("potato", "------------------------------------------------------------", t3)
	file2 := r.WriteFile("potato", "smaller but same date", t3)
	fstest.CheckItems(t, r.Fremote, file1)
	fstest.CheckItems(t, r.Flocal, file2)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Flocal, file2)
	fstest.CheckItems(t, r.Fremote, file2)
}

// Sync after changing a file's contents, changing modtime but length
// remaining the same
func TestSyncAfterChangingContentsOnly(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	var file1 fstest.Item
	if r.Fremote.Precision() == fs.ModTimeNotSupported {
		t.Logf("ModTimeNotSupported so forcing file to be a different size")
		file1 = r.WriteObject("potato", "different size to make sure it syncs", t3)
	} else {
		file1 = r.WriteObject("potato", "smaller but same date", t3)
	}
	file2 := r.WriteFile("potato", "SMALLER BUT SAME DATE", t2)
	fstest.CheckItems(t, r.Fremote, file1)
	fstest.CheckItems(t, r.Flocal, file2)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Flocal, file2)
	fstest.CheckItems(t, r.Fremote, file2)
}

// Sync after removing a file and adding a file --dry-run
func TestSyncAfterRemovingAFileAndAddingAFileDryRun(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteFile("potato2", "------------------------------------------------------------", t1)
	file2 := r.WriteObject("potato", "SMALLER BUT SAME DATE", t2)
	file3 := r.WriteBoth("empty space", "", t2)

	fs.Config.DryRun = true
	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	fs.Config.DryRun = false
	require.NoError(t, err)

	fstest.CheckItems(t, r.Flocal, file3, file1)
	fstest.CheckItems(t, r.Fremote, file3, file2)
}

// Sync after removing a file and adding a file
func TestSyncAfterRemovingAFileAndAddingAFile(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteFile("potato2", "------------------------------------------------------------", t1)
	file2 := r.WriteObject("potato", "SMALLER BUT SAME DATE", t2)
	file3 := r.WriteBoth("empty space", "", t2)
	fstest.CheckItems(t, r.Fremote, file2, file3)
	fstest.CheckItems(t, r.Flocal, file1, file3)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Flocal, file1, file3)
	fstest.CheckItems(t, r.Fremote, file1, file3)
}

// Sync after removing a file and adding a file
func TestSyncAfterRemovingAFileAndAddingAFileSubDir(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteFile("a/potato2", "------------------------------------------------------------", t1)
	file2 := r.WriteObject("b/potato", "SMALLER BUT SAME DATE", t2)
	file3 := r.WriteBoth("c/non empty space", "AhHa!", t2)
	require.NoError(t, fs.Mkdir(r.Fremote, "d"))
	require.NoError(t, fs.Mkdir(r.Fremote, "d/e"))

	fstest.CheckListingWithPrecision(
		t,
		r.Flocal,
		[]fstest.Item{
			file1,
			file3,
		},
		[]string{
			"a",
			"c",
		},
		fs.Config.ModifyWindow,
	)
	fstest.CheckListingWithPrecision(
		t,
		r.Fremote,
		[]fstest.Item{
			file2,
			file3,
		},
		[]string{
			"b",
			"c",
			"d",
			"d/e",
		},
		fs.Config.ModifyWindow,
	)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckListingWithPrecision(
		t,
		r.Flocal,
		[]fstest.Item{
			file1,
			file3,
		},
		[]string{
			"a",
			"c",
		},
		fs.Config.ModifyWindow,
	)
	fstest.CheckListingWithPrecision(
		t,
		r.Fremote,
		[]fstest.Item{
			file1,
			file3,
		},
		[]string{
			"a",
			"c",
		},
		fs.Config.ModifyWindow,
	)
}

// Sync after removing a file and adding a file with IO Errors
func TestSyncAfterRemovingAFileAndAddingAFileSubDirWithErrors(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteFile("a/potato2", "------------------------------------------------------------", t1)
	file2 := r.WriteObject("b/potato", "SMALLER BUT SAME DATE", t2)
	file3 := r.WriteBoth("c/non empty space", "AhHa!", t2)
	require.NoError(t, fs.Mkdir(r.Fremote, "d"))

	fstest.CheckListingWithPrecision(
		t,
		r.Flocal,
		[]fstest.Item{
			file1,
			file3,
		},
		[]string{
			"a",
			"c",
		},
		fs.Config.ModifyWindow,
	)
	fstest.CheckListingWithPrecision(
		t,
		r.Fremote,
		[]fstest.Item{
			file2,
			file3,
		},
		[]string{
			"b",
			"c",
			"d",
		},
		fs.Config.ModifyWindow,
	)

	fs.Stats.ResetCounters()
	fs.Stats.Error(nil)
	err := fs.Sync(r.Fremote, r.Flocal)
	assert.Equal(t, fs.ErrorNotDeleting, err)

	fstest.CheckListingWithPrecision(
		t,
		r.Flocal,
		[]fstest.Item{
			file1,
			file3,
		},
		[]string{
			"a",
			"c",
		},
		fs.Config.ModifyWindow,
	)
	fstest.CheckListingWithPrecision(
		t,
		r.Fremote,
		[]fstest.Item{
			file1,
			file2,
			file3,
		},
		[]string{
			"a",
			"b",
			"c",
			"d",
		},
		fs.Config.ModifyWindow,
	)
}

// Sync test delete after
func TestSyncDeleteAfter(t *testing.T) {
	// This is the default so we've checked this already
	// check it is the default
	require.Equal(t, fs.Config.DeleteMode, fs.DeleteModeAfter, "Didn't default to --delete-after")
}

// Sync test delete during
func TestSyncDeleteDuring(t *testing.T) {
	fs.Config.DeleteMode = fs.DeleteModeDuring
	defer func() {
		fs.Config.DeleteMode = fs.DeleteModeDefault
	}()

	TestSyncAfterRemovingAFileAndAddingAFile(t)
}

// Sync test delete before
func TestSyncDeleteBefore(t *testing.T) {
	fs.Config.DeleteMode = fs.DeleteModeBefore
	defer func() {
		fs.Config.DeleteMode = fs.DeleteModeDefault
	}()

	TestSyncAfterRemovingAFileAndAddingAFile(t)
}

// Copy test delete before - shouldn't delete anything
func TestCopyDeleteBefore(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()

	fs.Config.DeleteMode = fs.DeleteModeBefore
	defer func() {
		fs.Config.DeleteMode = fs.DeleteModeDefault
	}()

	file1 := r.WriteObject("potato", "hopefully not deleted", t1)
	file2 := r.WriteFile("potato2", "hopefully copied in", t1)
	fstest.CheckItems(t, r.Fremote, file1)
	fstest.CheckItems(t, r.Flocal, file2)

	fs.Stats.ResetCounters()
	err := fs.CopyDir(r.Fremote, r.Flocal)
	require.NoError(t, err)

	fstest.CheckItems(t, r.Fremote, file1, file2)
	fstest.CheckItems(t, r.Flocal, file2)
}

// Test with exclude
func TestSyncWithExclude(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteBoth("potato2", "------------------------------------------------------------", t1)
	file2 := r.WriteBoth("empty space", "", t2)
	file3 := r.WriteFile("enormous", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", t1) // 100 bytes
	fstest.CheckItems(t, r.Fremote, file1, file2)
	fstest.CheckItems(t, r.Flocal, file1, file2, file3)

	fs.Config.Filter.MaxSize = 40
	defer func() {
		fs.Config.Filter.MaxSize = -1
	}()

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Fremote, file2, file1)

	// Now sync the other way round and check enormous doesn't get
	// deleted as it is excluded from the sync
	fs.Stats.ResetCounters()
	err = fs.Sync(r.Flocal, r.Fremote)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Flocal, file2, file1, file3)
}

// Test with exclude and delete excluded
func TestSyncWithExcludeAndDeleteExcluded(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	file1 := r.WriteBoth("potato2", "------------------------------------------------------------", t1) // 60 bytes
	file2 := r.WriteBoth("empty space", "", t2)
	file3 := r.WriteBoth("enormous", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", t1) // 100 bytes
	fstest.CheckItems(t, r.Fremote, file1, file2, file3)
	fstest.CheckItems(t, r.Flocal, file1, file2, file3)

	fs.Config.Filter.MaxSize = 40
	fs.Config.Filter.DeleteExcluded = true
	defer func() {
		fs.Config.Filter.MaxSize = -1
		fs.Config.Filter.DeleteExcluded = false
	}()

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Fremote, file2)

	// Check sync the other way round to make sure enormous gets
	// deleted even though it is excluded
	fs.Stats.ResetCounters()
	err = fs.Sync(r.Flocal, r.Fremote)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Flocal, file2)
}

// Test with UpdateOlder set
func TestSyncWithUpdateOlder(t *testing.T) {
	if fs.Config.ModifyWindow == fs.ModTimeNotSupported {
		t.Skip("Can't run this test on fs which doesn't support mod time")
	}
	r := fstest.NewRun(t)
	defer r.Finalise()
	t2plus := t2.Add(time.Second / 2)
	t2minus := t2.Add(time.Second / 2)
	oneF := r.WriteFile("one", "one", t1)
	twoF := r.WriteFile("two", "two", t3)
	threeF := r.WriteFile("three", "three", t2)
	fourF := r.WriteFile("four", "four", t2)
	fiveF := r.WriteFile("five", "five", t2)
	fstest.CheckItems(t, r.Flocal, oneF, twoF, threeF, fourF, fiveF)
	oneO := r.WriteObject("one", "ONE", t2)
	twoO := r.WriteObject("two", "TWO", t2)
	threeO := r.WriteObject("three", "THREE", t2plus)
	fourO := r.WriteObject("four", "FOURFOUR", t2minus)
	fstest.CheckItems(t, r.Fremote, oneO, twoO, threeO, fourO)

	fs.Config.UpdateOlder = true
	oldModifyWindow := fs.Config.ModifyWindow
	fs.Config.ModifyWindow = fs.ModTimeNotSupported
	defer func() {
		fs.Config.UpdateOlder = false
		fs.Config.ModifyWindow = oldModifyWindow
	}()

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Fremote, oneO, twoF, threeO, fourF, fiveF)
}

// Test with TrackRenames set
func TestSyncWithTrackRenames(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()

	fs.Config.TrackRenames = true
	defer func() {
		fs.Config.TrackRenames = false

	}()

	haveHash := r.Fremote.Hashes().Overlap(r.Flocal.Hashes()).GetOne() != fs.HashNone
	canTrackRenames := haveHash && fs.CanServerSideMove(r.Fremote)
	t.Logf("Can track renames: %v", canTrackRenames)

	f1 := r.WriteFile("potato", "Potato Content", t1)
	f2 := r.WriteFile("yam", "Yam Content", t2)

	fs.Stats.ResetCounters()
	require.NoError(t, fs.Sync(r.Fremote, r.Flocal))

	fstest.CheckItems(t, r.Fremote, f1, f2)
	fstest.CheckItems(t, r.Flocal, f1, f2)

	// Now rename locally.
	f2 = r.RenameFile(f2, "yaml")

	fs.Stats.ResetCounters()
	require.NoError(t, fs.Sync(r.Fremote, r.Flocal))

	fstest.CheckItems(t, r.Fremote, f1, f2)

	if canTrackRenames {
		assert.Equal(t, fs.Stats.GetTransfers(), int64(0))
	} else {
		assert.Equal(t, fs.Stats.GetTransfers(), int64(1))
	}
}

// Test a server side move if possible, or the backup path if not
func testServerSideMove(t *testing.T, r *fstest.Run, withFilter bool) {
	FremoteMove, _, finaliseMove, err := fstest.RandomRemote(*fstest.RemoteName, *fstest.SubDir)
	require.NoError(t, err)
	defer finaliseMove()

	file1 := r.WriteBoth("potato2", "------------------------------------------------------------", t1)
	file2 := r.WriteBoth("empty space", "", t2)
	file3u := r.WriteBoth("potato3", "------------------------------------------------------------ UPDATED", t2)

	fstest.CheckItems(t, r.Fremote, file2, file1, file3u)

	t.Logf("Server side move (if possible) %v -> %v", r.Fremote, FremoteMove)

	// Write just one file in the new remote
	r.WriteObjectTo(FremoteMove, "empty space", "", t2, false)
	file3 := r.WriteObjectTo(FremoteMove, "potato3", "------------------------------------------------------------", t1, false)
	fstest.CheckItems(t, FremoteMove, file2, file3)

	// Do server side move
	fs.Stats.ResetCounters()
	err = fs.MoveDir(FremoteMove, r.Fremote)
	require.NoError(t, err)

	if withFilter {
		fstest.CheckItems(t, r.Fremote, file2)
	} else {
		fstest.CheckItems(t, r.Fremote)
	}
	fstest.CheckItems(t, FremoteMove, file2, file1, file3u)

	// Create a new empty remote for stuff to be moved into
	FremoteMove2, _, finaliseMove2, err := fstest.RandomRemote(*fstest.RemoteName, *fstest.SubDir)
	require.NoError(t, err)
	defer finaliseMove2()

	// Move it back to a new empty remote, dst does not exist this time
	fs.Stats.ResetCounters()
	err = fs.MoveDir(FremoteMove2, FremoteMove)
	require.NoError(t, err)

	if withFilter {
		fstest.CheckItems(t, FremoteMove2, file1, file3u)
		fstest.CheckItems(t, FremoteMove, file2)
	} else {
		fstest.CheckItems(t, FremoteMove2, file2, file1, file3u)
		fstest.CheckItems(t, FremoteMove)
	}
}

// Test a server side move if possible, or the backup path if not
func TestServerSideMove(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()
	testServerSideMove(t, r, false)
}

// Test a server side move if possible, or the backup path if not
func TestServerSideMoveWithFilter(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()

	fs.Config.Filter.MinSize = 40
	defer func() {
		fs.Config.Filter.MinSize = -1
	}()

	testServerSideMove(t, r, true)
}

// Test a server side move with overlap
func TestServerSideMoveOverlap(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()

	if r.Fremote.Features().DirMove != nil {
		t.Skip("Skipping test as remote supports DirMove")
	}

	subRemoteName := r.FremoteName + "/rclone-move-test"
	FremoteMove, err := fs.NewFs(subRemoteName)
	require.NoError(t, err)

	file1 := r.WriteObject("potato2", "------------------------------------------------------------", t1)
	fstest.CheckItems(t, r.Fremote, file1)

	// Subdir move with no filters should return ErrorCantMoveOverlapping
	err = fs.MoveDir(FremoteMove, r.Fremote)
	assert.EqualError(t, err, fs.ErrorCantMoveOverlapping.Error())

	// Now try with a filter which should also fail with ErrorCantMoveOverlapping
	fs.Config.Filter.MinSize = 40
	defer func() {
		fs.Config.Filter.MinSize = -1
	}()
	err = fs.MoveDir(FremoteMove, r.Fremote)
	assert.EqualError(t, err, fs.ErrorCantMoveOverlapping.Error())
}

// Test with BackupDir set
func testSyncBackupDir(t *testing.T, suffix string) {
	r := fstest.NewRun(t)
	defer r.Finalise()

	if !fs.CanServerSideMove(r.Fremote) {
		t.Skip("Skipping test as remote does not support server side move")
	}
	r.Mkdir(r.Fremote)

	fs.Config.BackupDir = r.FremoteName + "/backup"
	fs.Config.Suffix = suffix
	defer func() {
		fs.Config.BackupDir = ""
		fs.Config.Suffix = ""
	}()

	// Make the setup so we have one, two, three in the dest
	// and one (different), two (same) in the source
	file1 := r.WriteObject("dst/one", "one", t1)
	file2 := r.WriteObject("dst/two", "two", t1)
	file3 := r.WriteObject("dst/three", "three", t1)
	file2a := r.WriteFile("two", "two", t1)
	file1a := r.WriteFile("one", "oneA", t2)

	fstest.CheckItems(t, r.Fremote, file1, file2, file3)
	fstest.CheckItems(t, r.Flocal, file1a, file2a)

	fdst, err := fs.NewFs(r.FremoteName + "/dst")
	require.NoError(t, err)

	fs.Stats.ResetCounters()
	err = fs.Sync(fdst, r.Flocal)
	require.NoError(t, err)

	// one should be moved to the backup dir and the new one installed
	file1.Path = "backup/one" + suffix
	file1a.Path = "dst/one"
	// two should be unchanged
	// three should be moved to the backup dir
	file3.Path = "backup/three" + suffix

	fstest.CheckItems(t, r.Fremote, file1, file2, file3, file1a)

	// Now check what happens if we do it again
	// Restore a different three and update one in the source
	file3a := r.WriteObject("dst/three", "threeA", t2)
	file1b := r.WriteFile("one", "oneBB", t3)
	fstest.CheckItems(t, r.Fremote, file1, file2, file3, file1a, file3a)

	// This should delete three and overwrite one again, checking
	// the files got overwritten correctly in backup-dir
	fs.Stats.ResetCounters()
	err = fs.Sync(fdst, r.Flocal)
	require.NoError(t, err)

	// one should be moved to the backup dir and the new one installed
	file1a.Path = "backup/one" + suffix
	file1b.Path = "dst/one"
	// two should be unchanged
	// three should be moved to the backup dir
	file3a.Path = "backup/three" + suffix

	fstest.CheckItems(t, r.Fremote, file1b, file2, file3a, file1a)
}
func TestSyncBackupDir(t *testing.T)           { testSyncBackupDir(t, "") }
func TestSyncBackupDirWithSuffix(t *testing.T) { testSyncBackupDir(t, ".bak") }

// Check we can sync two files with differing UTF-8 representations
func TestSyncUTFNorm(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("Can't test UTF normalization on OS X")
	}

	r := fstest.NewRun(t)
	defer r.Finalise()

	// Two strings with different unicode normalization (from OS X)
	Encoding1 := "Testêé"
	Encoding2 := "Testêé"
	assert.NotEqual(t, Encoding1, Encoding2)
	assert.Equal(t, norm.NFC.String(Encoding1), norm.NFC.String(Encoding2))

	file1 := r.WriteFile(Encoding1, "This is a test", t1)
	fstest.CheckItems(t, r.Flocal, file1)

	file2 := r.WriteObject(Encoding2, "This is a old test", t2)
	fstest.CheckItems(t, r.Fremote, file2)

	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)

	// We should have transferred exactly one file, but kept the
	// normalized state of the file.
	assert.Equal(t, int64(1), fs.Stats.GetTransfers())
	fstest.CheckItems(t, r.Flocal, file1)
	file1.Path = file2.Path
	fstest.CheckItems(t, r.Fremote, file1)
}

// Test --immutable
func TestSyncImmutable(t *testing.T) {
	r := fstest.NewRun(t)
	defer r.Finalise()

	fs.Config.Immutable = true
	defer func() { fs.Config.Immutable = false }()

	// Create file on source
	file1 := r.WriteFile("existing", "potato", t1)
	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote)

	// Should succeed
	fs.Stats.ResetCounters()
	err := fs.Sync(r.Fremote, r.Flocal)
	require.NoError(t, err)
	fstest.CheckItems(t, r.Flocal, file1)
	fstest.CheckItems(t, r.Fremote, file1)

	// Modify file data and timestamp on source
	file2 := r.WriteFile("existing", "tomatoes", t2)
	fstest.CheckItems(t, r.Flocal, file2)
	fstest.CheckItems(t, r.Fremote, file1)

	// Should fail with ErrorImmutableModified and not modify local or remote files
	fs.Stats.ResetCounters()
	err = fs.Sync(r.Fremote, r.Flocal)
	assert.EqualError(t, err, fs.ErrorImmutableModified.Error())
	fstest.CheckItems(t, r.Flocal, file2)
	fstest.CheckItems(t, r.Fremote, file1)
}
