/*

Carta is a simple SQL data mapper for complex Go structs.
Loads SQL data onto Go structs while keeping track of has-one and has-many relationships

To use carta:

1) Run your query
if rows, err = sqlDB.Query(blogQuery); err != nil {
	// error
}

2) Instantiate a slice(or struct) which you want to populate
blogs := []Blog{}

3) Map the SQL rows to your slice
if err := carta.Map(rows, &blogs); err != nil {
	// error
}

For more examples and guides go to: https://jackskj.github.io/carta/

*/

package carta
