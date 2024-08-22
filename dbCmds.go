package main


import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)



func getAllTasks(db *sql.DB) ([]task, error){

    rows, err := db.Query("SELECT * FROM tasks")
    defer rows.Close()

    tasks := make([]task, 0)

    for rows.Next() {
        newTask := task{}
        err = rows.Scan(&newTask.id, &newTask.title, &newTask.desc, &newTask.prio, &newTask.currentStatus)

        tasks = append(tasks, newTask)
    }
    return tasks, err
}


func insertNewTask(db *sql.DB, newTask task) error {

    stmt, err := db.Prepare("INSERT INTO tasks (title, desc, prio, status) VALUES (?, ?, ?, ?)")
    stmt.Exec(newTask.title, newTask.desc, 1, todo)
    defer stmt.Close()
    return err
}
