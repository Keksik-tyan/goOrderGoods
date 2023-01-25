package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type product struct {
	product_id    int
	product_name  string
	order_product order_product
	product_box   product_box
	box           box
}

type order_product struct {
	order_id int
	quantity int
}

type box struct {
	box_name string
}

type product_box struct {
	is_main bool
}

func main() {

	consoleNumberOrder := os.Args[1:]

	var NumberOrder string

	if len(consoleNumberOrder) != 0 {
		NumberOrder = strings.Join(consoleNumberOrder, ", ")
	} else {
		panic("Пустые аргументы в программе")
	}

	fmt.Println("Страница сборки заказов - ", NumberOrder)

	conn := "user=postgres password=1234 dbname=new_order_goods sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sql := fmt.Sprintf(`select distinct
    p.product_id,
    p.product_name,
    op.order_id,
    op.quantity,
    pb.is_main,
    b.box_name  
    from product p 
    left join order_product op on p.product_id = op.product_id
    left join product_box pb on p.product_id = pb.product_id
    left join box b on pb.box_id = b.box_id
    where op.order_id in (%s)
    order by  p.product_id, pb.is_main, b.box_name`, NumberOrder)

	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	products := []product{}

	for rows.Next() {
		p := product{}
		err := rows.Scan(
			&p.product_id,
			&p.product_name,
			&p.order_product.order_id,
			&p.order_product.quantity,
			&p.product_box.is_main,
			&p.box.box_name)

		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	subs := getSubs(products)
	boxName := ""
	for _, p := range products {
		if p.product_box.is_main {
			if p.box.box_name != boxName {
				fmt.Println("---Стеллаж", p.box.box_name, "---")
				boxName = p.box.box_name
			}
			fmt.Println(p.product_name, "(id =", p.product_id, ")")
			fmt.Println("заказ", p.order_product.order_id, ",", p.order_product.quantity, "шт")

			sub := subs[p.product_id]
			if len(sub) != 0 {
				fmt.Println("доп.стеллаж: ", strings.Join(sub, ", "))
			}

		}
	}
}

func getSubs(p []product) map[int][]string {
	res := make(map[int][]string)
	for _, pr := range p {
		if !pr.product_box.is_main {
			arr := res[pr.product_id]
			res[pr.product_id] = append(arr, pr.box.box_name)
		}
	}
	return res
}
