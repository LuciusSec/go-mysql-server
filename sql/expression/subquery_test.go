package expression_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xsec-lab/go-mysql-server/memory"
	"github.com/xsec-lab/go-mysql-server/sql"
	"github.com/xsec-lab/go-mysql-server/sql/expression"
	"github.com/xsec-lab/go-mysql-server/sql/plan"
)

func TestSubquery(t *testing.T) {
	require := require.New(t)
	table := memory.NewTable("", nil)
	require.NoError(table.Insert(sql.NewEmptyContext(), nil))

	subquery := expression.NewSubquery(plan.NewProject(
		[]sql.Expression{
			expression.NewLiteral("one", sql.Text),
		},
		plan.NewResolvedTable(table),
	))

	value, err := subquery.Eval(sql.NewEmptyContext(), nil)
	require.NoError(err)
	require.Equal(value, "one")
}

func TestSubqueryTooManyRows(t *testing.T) {
	require := require.New(t)
	table := memory.NewTable("", nil)
	require.NoError(table.Insert(sql.NewEmptyContext(), nil))
	require.NoError(table.Insert(sql.NewEmptyContext(), nil))

	subquery := expression.NewSubquery(plan.NewProject(
		[]sql.Expression{
			expression.NewLiteral("one", sql.Text),
		},
		plan.NewResolvedTable(table),
	))

	_, err := subquery.Eval(sql.NewEmptyContext(), nil)
	require.Error(err)
}

func TestSubqueryMultipleRows(t *testing.T) {
	require := require.New(t)

	ctx := sql.NewEmptyContext()
	table := memory.NewTable("foo", sql.Schema{
		{Name: "t", Source: "foo", Type: sql.Text},
	})

	require.NoError(table.Insert(ctx, sql.Row{"one"}))
	require.NoError(table.Insert(ctx, sql.Row{"two"}))
	require.NoError(table.Insert(ctx, sql.Row{"three"}))

	subquery := expression.NewSubquery(plan.NewProject(
		[]sql.Expression{
			expression.NewGetField(0, sql.Text, "t", false),
		},
		plan.NewResolvedTable(table),
	))

	values, err := subquery.EvalMultiple(ctx)
	require.NoError(err)
	require.Equal(values, []interface{}{"one", "two", "three"})
}
