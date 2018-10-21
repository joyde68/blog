package routes


// CommentHtml returns rendered comment template html with own content.
/*
func CommentHtml(context *macaron.Context, c *models.Content) (string, error) {
	thm := Theme(context)
	if !thm.Has("comment.html") {
		return ""
	}
	return thm.Tpl("comment", map[string]interface{}{
		"Content":  c,
		"Comments": c.Comments,
	})
}
*/

// SidebarHtml returns rendered sidebar template html.
/*
func SidebarHtml(context *macaron.Context) string {
	thm := Theme(context)
	if !thm.Has("sidebar.html") {
		return ""
	}
	popSize, _ := strconv.Atoi(models.GetSetting("popular_size"))
	if popSize < 1 {
		popSize = 4
	}
	cmtSize, _ := strconv.Atoi(models.GetSetting("recent_comment_size"))
	if cmtSize < 1 {
		cmtSize = 3
	}
	return thm.Tpl("sidebar", map[string]interface{}{
		"Popular":       models.GetPopularArticleList(popSize),
		"RecentComment": models.GetCommentRecentList(cmtSize),
		"Tags":          models.GetContentTags(),
	})
}
*/