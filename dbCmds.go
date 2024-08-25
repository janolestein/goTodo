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
func deleteTask(db *sql.DB, id int) error {

    stmt, err := db.Prepare("DELETE FROM tasks WHERE task_id = (?)")
    stmt.Exec(id)
    defer stmt.Close()
    return err
}

func updateTask(db *sql.DB, taskToUpdate task) error {

    stmt, err := db.Prepare("UPDATE tasks SET title = ?, desc = ?, status = ? WHERE task_id = ?")
    stmt.Exec(taskToUpdate.title, taskToUpdate.desc, taskToUpdate.currentStatus, taskToUpdate.id)
    defer stmt.Close()
    return err
}
