package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
)

// loads the session state at the time that a custom command is invoked, for use
// in the custom command's template strings
type SessionStateLoader struct {
	c          *helpers.HelperCommon
	refsHelper *helpers.RefsHelper
}

func NewSessionStateLoader(c *helpers.HelperCommon, refsHelper *helpers.RefsHelper) *SessionStateLoader {
	return &SessionStateLoader{
		c:          c,
		refsHelper: refsHelper,
	}
}

// We create shims for all the model classes in order to get a more stable API
// for custom commands. At the moment these are almost identical to the model
// classes, but this allows us to add "private" fields to the model classes that
// we don't want to expose to custom commands, or rename a model field to a
// better name without breaking people's custom commands. In such a case we add
// the new, better name to the shim but keep the old one for backwards
// compatibility. We already did this for Commit.Sha, which was renamed to Hash.

type CommitShim struct {
	Hash          string // deprecated: use Sha
	Sha           string
	Name          string
	Status        models.CommitStatus
	Action        todo.TodoCommand
	Tags          []string
	ExtraInfo     string
	AuthorName    string
	AuthorEmail   string
	UnixTimestamp int64
	Divergence    models.Divergence
	Parents       []string
}

func commitShimFromModelCommit(commit *models.Commit) *CommitShim {
	if commit == nil {
		return nil
	}

	return &CommitShim{
		Hash:          commit.Hash,
		Sha:           commit.Hash,
		Name:          commit.Name,
		Status:        commit.Status,
		Action:        commit.Action,
		Tags:          commit.Tags,
		ExtraInfo:     commit.ExtraInfo,
		AuthorName:    commit.AuthorName,
		AuthorEmail:   commit.AuthorEmail,
		UnixTimestamp: commit.UnixTimestamp,
		Divergence:    commit.Divergence,
		Parents:       commit.Parents,
	}
}

type FileShim struct {
	Name                    string
	PreviousName            string
	HasStagedChanges        bool
	HasUnstagedChanges      bool
	Tracked                 bool
	Added                   bool
	Deleted                 bool
	HasMergeConflicts       bool
	HasInlineMergeConflicts bool
	DisplayString           string
	ShortStatus             string
	IsWorktree              bool
}

func fileShimFromModelFile(file *models.File) *FileShim {
	if file == nil {
		return nil
	}

	return &FileShim{
		Name:                    file.Name,
		PreviousName:            file.PreviousName,
		HasStagedChanges:        file.HasStagedChanges,
		HasUnstagedChanges:      file.HasUnstagedChanges,
		Tracked:                 file.Tracked,
		Added:                   file.Added,
		Deleted:                 file.Deleted,
		HasMergeConflicts:       file.HasMergeConflicts,
		HasInlineMergeConflicts: file.HasInlineMergeConflicts,
		DisplayString:           file.DisplayString,
		ShortStatus:             file.ShortStatus,
		IsWorktree:              file.IsWorktree,
	}
}

type BranchShim struct {
	Name           string
	DisplayName    string
	Recency        string
	Pushables      string // deprecated: use AheadForPull
	Pullables      string // deprecated: use BehindForPull
	AheadForPull   string
	BehindForPull  string
	AheadForPush   string
	BehindForPush  string
	UpstreamGone   bool
	Head           bool
	DetachedHead   bool
	UpstreamRemote string
	UpstreamBranch string
	Subject        string
	CommitHash     string
}

func branchShimFromModelBranch(branch *models.Branch) *BranchShim {
	if branch == nil {
		return nil
	}

	return &BranchShim{
		Name:           branch.Name,
		DisplayName:    branch.DisplayName,
		Recency:        branch.Recency,
		Pushables:      branch.AheadForPull,
		Pullables:      branch.BehindForPull,
		AheadForPull:   branch.AheadForPull,
		BehindForPull:  branch.BehindForPull,
		AheadForPush:   branch.AheadForPush,
		BehindForPush:  branch.BehindForPush,
		UpstreamGone:   branch.UpstreamGone,
		Head:           branch.Head,
		DetachedHead:   branch.DetachedHead,
		UpstreamRemote: branch.UpstreamRemote,
		UpstreamBranch: branch.UpstreamBranch,
		Subject:        branch.Subject,
		CommitHash:     branch.CommitHash,
	}
}

type RemoteBranchShim struct {
	Name       string
	RemoteName string
}

func remoteBranchShimFromModelRemoteBranch(remoteBranch *models.RemoteBranch) *RemoteBranchShim {
	if remoteBranch == nil {
		return nil
	}

	return &RemoteBranchShim{
		Name:       remoteBranch.Name,
		RemoteName: remoteBranch.RemoteName,
	}
}

type RemoteShim struct {
	Name     string
	Urls     []string
	Branches []*RemoteBranchShim
}

func remoteShimFromModelRemote(remote *models.Remote) *RemoteShim {
	if remote == nil {
		return nil
	}

	return &RemoteShim{
		Name: remote.Name,
		Urls: remote.Urls,
		Branches: lo.Map(remote.Branches, func(branch *models.RemoteBranch, _ int) *RemoteBranchShim {
			return remoteBranchShimFromModelRemoteBranch(branch)
		}),
	}
}

type TagShim struct {
	Name    string
	Message string
}

func tagShimFromModelRemote(tag *models.Tag) *TagShim {
	if tag == nil {
		return nil
	}

	return &TagShim{
		Name:    tag.Name,
		Message: tag.Message,
	}
}

type StashEntryShim struct {
	Index   int
	Recency string
	Name    string
}

func stashEntryShimFromModelRemote(stashEntry *models.StashEntry) *StashEntryShim {
	if stashEntry == nil {
		return nil
	}

	return &StashEntryShim{
		Index:   stashEntry.Index,
		Recency: stashEntry.Recency,
		Name:    stashEntry.Name,
	}
}

type CommitFileShim struct {
	Name         string
	ChangeStatus string
}

func commitFileShimFromModelRemote(commitFile *models.CommitFile) *CommitFileShim {
	if commitFile == nil {
		return nil
	}

	return &CommitFileShim{
		Name:         commitFile.Name,
		ChangeStatus: commitFile.ChangeStatus,
	}
}

type WorktreeShim struct {
	IsMain        bool
	IsCurrent     bool
	Path          string
	IsPathMissing bool
	GitDir        string
	Branch        string
	Name          string
}

func worktreeShimFromModelRemote(worktree *models.Worktree) *WorktreeShim {
	if worktree == nil {
		return nil
	}

	return &WorktreeShim{
		IsMain:        worktree.IsMain,
		IsCurrent:     worktree.IsCurrent,
		Path:          worktree.Path,
		IsPathMissing: worktree.IsPathMissing,
		GitDir:        worktree.GitDir,
		Branch:        worktree.Branch,
		Name:          worktree.Name,
	}
}

// SessionState captures the current state of the application for use in custom commands
type SessionState struct {
	SelectedLocalCommit    *CommitShim
	SelectedReflogCommit   *CommitShim
	SelectedSubCommit      *CommitShim
	SelectedFile           *FileShim
	SelectedPath           string
	SelectedLocalBranch    *BranchShim
	SelectedRemoteBranch   *RemoteBranchShim
	SelectedRemote         *RemoteShim
	SelectedTag            *TagShim
	SelectedStashEntry     *StashEntryShim
	SelectedCommitFile     *CommitFileShim
	SelectedCommitFilePath string
	SelectedWorktree       *WorktreeShim
	CheckedOutBranch       *BranchShim
}

func (self *SessionStateLoader) call() *SessionState {
	return &SessionState{
		SelectedFile:           fileShimFromModelFile(self.c.Contexts().Files.GetSelectedFile()),
		SelectedPath:           self.c.Contexts().Files.GetSelectedPath(),
		SelectedLocalCommit:    commitShimFromModelCommit(self.c.Contexts().LocalCommits.GetSelected()),
		SelectedReflogCommit:   commitShimFromModelCommit(self.c.Contexts().ReflogCommits.GetSelected()),
		SelectedLocalBranch:    branchShimFromModelBranch(self.c.Contexts().Branches.GetSelected()),
		SelectedRemoteBranch:   remoteBranchShimFromModelRemoteBranch(self.c.Contexts().RemoteBranches.GetSelected()),
		SelectedRemote:         remoteShimFromModelRemote(self.c.Contexts().Remotes.GetSelected()),
		SelectedTag:            tagShimFromModelRemote(self.c.Contexts().Tags.GetSelected()),
		SelectedStashEntry:     stashEntryShimFromModelRemote(self.c.Contexts().Stash.GetSelected()),
		SelectedCommitFile:     commitFileShimFromModelRemote(self.c.Contexts().CommitFiles.GetSelectedFile()),
		SelectedCommitFilePath: self.c.Contexts().CommitFiles.GetSelectedPath(),
		SelectedSubCommit:      commitShimFromModelCommit(self.c.Contexts().SubCommits.GetSelected()),
		SelectedWorktree:       worktreeShimFromModelRemote(self.c.Contexts().Worktrees.GetSelected()),
		CheckedOutBranch:       branchShimFromModelBranch(self.refsHelper.GetCheckedOutRef()),
	}
}
