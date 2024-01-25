package sess

import "xorm.io/xorm"

func Transaction(engine *xorm.Engine, f func(session *xorm.Session) error) error {
	session := engine.NewSession()
	if err := session.Begin(); err != nil {
		return err
	}
	defer session.Close()
	if err := f(session); err != nil {
		return err
	}
	return session.Commit()
}
