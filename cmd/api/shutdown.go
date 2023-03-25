package main

func (app *application) shutdown() {
	app.dbpool.Close()
}
