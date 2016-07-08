package main

import (
	"fmt"
)

func (q *Query) String() string {
	return fmt.Sprintf("%v -> %v", q.Statement, q.t)
}

func (db *Database) String() string {
	return fmt.Sprintf("DB %v", db.Name)
}

func (rp *RetentionPolicy) String() string {
	return fmt.Sprintf(`  RP %v -> %v
    Default -> %v`, rp.Name, rp.Duration, rp.Default)
}

func (m *Measurement) String() string {
	return fmt.Sprintf("  M %v -> %v", m.Name, m.Series)
}

func (t *Tag) String() string {
	return fmt.Sprintf("    T %v -> %v", t.Name, t.Cardinality)
}

func (f *Field) String() string {
	return fmt.Sprintf("    F %v -> %v", f.Name, f.Type)
}
