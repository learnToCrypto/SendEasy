package pagination

type Paginator struct {
	Limit int

	PageNum  int
	Current  int
	PageLink []int // page links around current
}

func NewPaginator(limit int, current int, pageNum int) *Paginator {
	p := Paginator{
		Limit:   limit,
		PageNum: pageNum,
		Current: current,
	}
	return &p
}

// Logic
func (p *Paginator) GeneratePageLink() {
	switch {
	case p.PageNum < 11:
		for i := 1; i <= p.PageNum; i++ {
			p.PageLink = append(p.PageLink, i)
		}
	case p.PageNum >= 11 && p.Current < 7:
		for i := 1; i <= 10; i++ {
			p.PageLink = append(p.PageLink, i)
		}
	case p.PageNum >= 11 && p.Current >= 7 && p.Current+4 > p.PageNum:
		for i := p.PageNum - 9; i <= p.PageNum; i++ {
			p.PageLink = append(p.PageLink, i)
		}
	default:
		for i := p.Current - 5; i <= p.Current+4; i++ {
			p.PageLink = append(p.PageLink, i)
		}
	}
}

//PageNUm calculates number of pages from Number of items in list and limit per page
func PageNum(itemNum, limit int) int {
	if (itemNum % limit) == 0 {
		return itemNum / limit
	} else {
		return itemNum/limit + 1
	}
}
