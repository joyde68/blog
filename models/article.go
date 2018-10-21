package models

import "github.com/joyde68/blog/pkg"

func generatePublishArticleIndex() {
	arr := make([]int, 0)
	for _, id := range contentsIndex["article"] {
		c := GetContentById(id)
		if c.Status == "publish" {
			arr = append(arr, id)
		}
	}
	contentsIndex["article-publish"] = arr
}

// GetPublishArticleList gets published article list and pager.
func GetPublishArticleList(page, size int) ([]*Content, *pkg.Pager) {
	index := contentsIndex["article-publish"]
	pager := pkg.NewPager(page, size, len(index))
	articles := make([]*Content, 0)
	if len(index) < 1 {
		return articles, pager
	}
	if page > pager.Pages {
		return articles, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		articles = append(articles, GetContentById(index[i-1]))
	}
	return articles, pager
}

// GetArticleList gets articles list and pager no matter article status.
func GetArticleList(page, size int) ([]*Content, *pkg.Pager) {
	index := contentsIndex["article"]
	pager := pkg.NewPager(page, size, len(index))
	articles := make([]*Content, 0)
	if len(index) < 1 {
		return articles, pager
	}
	if page > pager.Pages {
		return articles, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		articles = append(articles, GetContentById(index[i-1]))
	}
	return articles, pager
}

// GetPopularArticleList returns popular articles list.
// Popular articles are ordered by comment number.
func GetPopularArticleList(size int) []*Content {
	index := contentsIndex["article-pop"]
	pager := pkg.NewPager(1, size, len(index))
	articles := make([]*Content, 0)
	if len(index) < 1 {
		return articles
	}
	if 1 > pager.Pages {
		return articles
	}
	for i := pager.Begin; i <= pager.End; i++ {
		articles = append(articles, GetContentById(index[i-1]))
	}
	return articles
}

// GetTaggedArticleList returns tagged articles list.
// These articles contains same one tag.
func GetTaggedArticleList(tag string, page, size int) ([]*Content, *pkg.Pager) {
	index := contentsIndex["t-"+tag]
	pager := pkg.NewPager(page, size, len(index))
	articles := make([]*Content, 0)
	if len(index) < 1 {
		return articles, pager
	}
	if page > pager.Pages {
		return articles, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		articles = append(articles, GetContentById(index[i-1]))
	}
	return articles, pager
}
