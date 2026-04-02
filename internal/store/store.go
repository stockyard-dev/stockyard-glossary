package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{ db *sql.DB }
type Term struct { ID string `json:"id"`; Name string `json:"name"`; Definition string `json:"definition"`; Category string `json:"category,omitempty"`; Aliases string `json:"aliases,omitempty"`; CreatedAt string `json:"created_at"` }
func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil { return nil, err }
	db, err := sql.Open("sqlite", filepath.Join(d, "glossary.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil { return nil, err }
	db.Exec(`CREATE TABLE IF NOT EXISTS terms (id TEXT PRIMARY KEY, name TEXT NOT NULL, definition TEXT DEFAULT '', category TEXT DEFAULT '', aliases TEXT DEFAULT '', created_at TEXT DEFAULT (datetime('now')))`)
	return &DB{db: db}, nil
}
func (d *DB) Close() error { return d.db.Close() }
func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string { return time.Now().UTC().Format(time.RFC3339) }
func (d *DB) Create(t *Term) error { t.ID=genID();t.CreatedAt=now(); _, err := d.db.Exec(`INSERT INTO terms VALUES(?,?,?,?,?,?)`,t.ID,t.Name,t.Definition,t.Category,t.Aliases,t.CreatedAt); return err }
func (d *DB) Get(id string) *Term { var t Term; if d.db.QueryRow(`SELECT * FROM terms WHERE id=?`,id).Scan(&t.ID,&t.Name,&t.Definition,&t.Category,&t.Aliases,&t.CreatedAt)!=nil{return nil}; return &t }
func (d *DB) List() []Term { rows,_:=d.db.Query(`SELECT * FROM terms ORDER BY name`); if rows==nil{return nil}; defer rows.Close(); var o []Term; for rows.Next(){var t Term;rows.Scan(&t.ID,&t.Name,&t.Definition,&t.Category,&t.Aliases,&t.CreatedAt);o=append(o,t)}; return o }
func (d *DB) Update(id string, t *Term) error { _,err:=d.db.Exec(`UPDATE terms SET name=?,definition=?,category=?,aliases=? WHERE id=?`,t.Name,t.Definition,t.Category,t.Aliases,id); return err }
func (d *DB) Delete(id string) error { _,err:=d.db.Exec(`DELETE FROM terms WHERE id=?`,id); return err }
func (d *DB) Search(q string) []Term { s:="%"+q+"%"; rows,_:=d.db.Query(`SELECT * FROM terms WHERE name LIKE ? OR definition LIKE ? OR aliases LIKE ? ORDER BY name`,s,s,s); if rows==nil{return nil}; defer rows.Close(); var o []Term; for rows.Next(){var t Term;rows.Scan(&t.ID,&t.Name,&t.Definition,&t.Category,&t.Aliases,&t.CreatedAt);o=append(o,t)}; return o }
func (d *DB) Categories() []string { rows,_:=d.db.Query(`SELECT DISTINCT category FROM terms WHERE category!='' ORDER BY category`); if rows==nil{return nil}; defer rows.Close(); var o []string; for rows.Next(){var c string;rows.Scan(&c);o=append(o,c)}; return o }
func (d *DB) Count() int { var n int; d.db.QueryRow(`SELECT COUNT(*) FROM terms`).Scan(&n); return n }
