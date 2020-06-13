package testdata

import (
	"database/sql"
)

type RelationTest struct {
	Id                int              `db:"id" json:"id,omitempty"`
	BasicSubMap       []sql.NullInt64  `db:"basic_submap" json:"basic_submap,omitempty"`
	BasicSubMapPtr    *[]sql.NullInt64 `db:"basic_submap_ptr" json:"basic_submap_ptr,omitempty"`
	BasicSubMapPtrPtr *[]*int64        `db:"basic_submap_ptr_ptr" json:"basic_submap_ptr_ptr,omitempty"`

	SubMap       []SampleSubmap   `db:"submap" json:"submap,omitempty"`
	SubMapPtr    *[]SampleSubmap  `db:"submap_ptr" json:"submap_ptr,omitempty"`
	SubMapPtrPtr *[]*SampleSubmap `db:"submap_ptr_ptr" json:"submap_ptr_ptr,omitempty"`
}

type SampleSubmap struct {
	SubId int `db:"sample_submap" json:"sample_submap,omitempty"`
}

var RelationTestQuery = `
select * from relations order by id,  basic_submap,  basic_submap_ptr,  basic_submap_ptr_ptr, submap_sample_submap, submap_ptr_sample_submap, submap_ptr_ptr_sample_submap;
`
