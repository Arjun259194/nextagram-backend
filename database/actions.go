package database

func AddUser(name,email,password string) error {
  stmt,err := Db.Prepare("INSERT INTO user (name, email, password) VALUES (?,?,?)")
  if err != nil {
    return err
  }

  defer stmt.Close()

  _,err = stmt.Exec(name,email,password)

  if err != nil {
    return err
  }

  return nil
}
