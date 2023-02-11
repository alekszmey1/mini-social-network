package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"homework.31/pkg/entity"
)

const s = "mysql"
const root = "root:root@tcp(127.0.0.1:3307)/skillbox_31"

// Определение хранилища, и типы данных, хранящиеся там
type repository struct {
}
type u struct {
	id      int
	name    string
	age     int
	friends string
}

// Функция создания хранилища
func NewRepository() *repository {
	return &repository{}
}

// Метод добавления в хранилище новую сущность
func (r *repository) CreateUser(user *entity.User) (int64, error) {
	db, err := sql.Open(s, root)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	result, err := db.Exec("insert into skill31 (name, age) values (?, ?)",
		user.Name, user.Age)
	if err != nil {
		panic(err)
	}
	user.Id, err = result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return user.Id, nil
}

// метод добавления друзей
func (r *repository) MakeFriends(friends *entity.MakeFriends) (a int, b int, err error) {
	db := dbOpen(s, root)
	defer db.Close()

	idFirstFriends := stringToInt(friends.SourceId)
	idSecondFriends := stringToInt(friends.TargetId)
	friend1 := getFriendsID(idFirstFriends, db)
	friend2 := getFriendsID(idSecondFriends, db)
	friend1 = friend1 + " " + friends.TargetId
	friend2 = friend2 + " " + friends.SourceId
	a = updateFriends(db, friend1, idFirstFriends)
	b = updateFriends(db, friend2, idSecondFriends)
	return a, b, nil
}

// удаление пользователя
func (r *repository) DeleteUser(user *entity.DeleteUser) string {
	db := dbOpen(s, root)
	defer db.Close()
	//получили айди номер удаляемого объекта и перевели значение в инт
	userInt := stringToInt(user.TargetId)
	fmt.Printf("удаляем друзей у %v\n", userInt)
	//получили список друзей, у которых мы значимся в друзьях
	friendsDelete := getFriendsID(userInt, db)
	fmt.Printf("получен список друзей, а именно %s \n", friendsDelete)
	if friendsDelete != "" {
		fmt.Println("Если список друзей есть, то мы начинаем с ним работать")
		//переделал строку со списком друзей в срез строк друзей
		format1 := strings.Trim(friendsDelete, " ")
		f := strings.Split(format1, " ")
		fmt.Printf("переделал строку со списком друзей в срез строк друзей %s\n", f)

		//пробегаем по срезу друзей и если есть совпадение, то необходимо открыть список друзей у второго и уменьшить строку на удаляемого
		//value - айди удаляемого
		for i, value := range f {
			//перевели значение айди у друга удаляемого в инт
			fmt.Printf("взят %v друг с айди %s\n", i, value)
			v := stringToInt(value)
			//получили список друзей строкой у друга удаляемого
			friends := getFriendsID(v, db)
			fmt.Printf("получили список друзей строкой у друга удаляемого %s\n", friends)
			//переделал список друзей в срез
			format := strings.Trim(friends, " ")
			f2 := strings.Split(format, " ")
			for i, value := range f2 {
				if value == user.TargetId {
					f2 = anotherRemove(f2, i)
					f3 := strings.Join(f2, " ")
					result, err := db.Exec(fmt.Sprintf("update skill31 set friends = '%s' where id = '%v'", f3, v))
					if err != nil {
						fmt.Println(err)
						panic(err)
					}
					fmt.Println(result.LastInsertId())
				}
			}
		}
		fmt.Println("удалили у друзей в их списках")
	}

	rows, err := db.Query(fmt.Sprintf("SELECT `name` FROM `skill31` where `id` = '%v'", userInt))
	fmt.Println("ВЫТАСКИВАЕМ имя удаляемого")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var user1 u
	for rows.Next() {
		err := rows.Scan(&user1.name)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	fmt.Printf("вытащили имя удаляемого %s", user1.name)
	fmt.Println("начинаем удалять из таблицы ")
	result, idFriends := db.Exec(fmt.Sprintf("delete from skill31 where id = '%s'", user.TargetId))
	if idFriends != nil {
		panic(idFriends)
	}
	fmt.Println(result.LastInsertId()) // id последнего удаленого объекта
	return user1.name
}

// получение списка друзей
func (r *repository) GetFriends(a int) (b string, err error) {
	db := dbOpen(s, root)
	b = getFriendsID(a, db)
	return b, nil
}

// обновление возраста
func (r *repository) UpdateAge(user *entity.UpdateUser) string {
	db := dbOpen(s, root)
	defer db.Close()
	idUpdate := stringToInt(user.Target)
	fmt.Println(idUpdate)
	newAge := stringToInt(user.NewAge)
	fmt.Println(newAge)
	result, err := db.Exec("update skill31 set age = ? where id = ?", newAge, idUpdate)
	if err != nil {

		panic(err)
	}
	fmt.Println(result.LastInsertId())
	return "возраст  пользователя успешно обновлен"
}

func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("err")
	}
	return i
}
func getFriendsID(i int, db *sql.DB) (s string) {
	log.Println("запустилась функция get Friends")
	log.Println(i)
	rows, _ := db.Query(fmt.Sprintf("SELECT `friends` FROM `skill31` where `id` = '%v'", i))
	defer rows.Close()
	var user u
	if rows == nil {
		return "список друзей пуст"
	} else {
		for rows.Next() {
			err := rows.Scan(&user.friends)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		return user.friends
	}
	/*if err != nil {
		panic(err)
	}*/
}

func updateFriends(db *sql.DB, s string, id int) int {
	result, err := db.Exec("update skill31 set friends = ? where id = ?", s, id)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.LastInsertId())
	return id
}
func dbOpen(str string, r string) (db *sql.DB) {
	db, err := sql.Open(str, r)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	return
}
func anotherRemove(s []string, i int) []string {
	s = append(s[:i], s[i+1:]...)
	return s
}
