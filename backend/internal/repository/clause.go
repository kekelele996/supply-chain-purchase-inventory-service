package repository

import "gorm.io/gorm/clause"

// clauseLocking 行级锁：SELECT ... FOR UPDATE，用于并发写场景。
var clauseLocking = clause.Locking{Strength: "UPDATE"}
