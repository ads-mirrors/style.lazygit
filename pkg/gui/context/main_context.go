package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MainContext struct {
	*SimpleContext
	*SearchTrait

	// The side panel context that owns this main context. Only set if the main
	// context is focused, not when the side panel is focused and is just
	// rendering to it.
	owningSidePanelContext types.Context
}

var _ types.ISearchableContext = (*MainContext)(nil)

func NewMainContext(
	view *gocui.View,
	windowName string,
	key types.ContextKey,
	c *ContextCommon,
) *MainContext {
	ctx := &MainContext{
		SimpleContext: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:             types.MAIN_CONTEXT,
				View:             view,
				WindowName:       windowName,
				Key:              key,
				Focusable:        true,
				HighlightOnFocus: false,
			})),
		SearchTrait: NewSearchTrait(c),
	}

	ctx.GetView().SetOnSelectItem(ctx.SearchTrait.onSelectItemWrapper(func(int) error { return nil }))

	return ctx
}

func (self *MainContext) ModelSearchResults(searchStr string, caseSensitive bool) []gocui.SearchPosition {
	return nil
}

func (self *MainContext) SetOwningSidePanelContext(context types.Context) {
	self.owningSidePanelContext = context
}

func (self *MainContext) GetOwningSidePanelContext() types.Context {
	return self.owningSidePanelContext
}
