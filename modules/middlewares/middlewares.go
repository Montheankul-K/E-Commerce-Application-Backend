// Middlewares : ตัวกลางระหว่าง user กับ api
package middlewares

type Role struct {
	Id    int    `db:"id"`
	Title string `db:"title"`
}
