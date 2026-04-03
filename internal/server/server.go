package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-glossary/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/terms",s.list);s.mux.HandleFunc("POST /api/terms",s.create);s.mux.HandleFunc("GET /api/terms/{id}",s.get);s.mux.HandleFunc("PUT /api/terms/{id}",s.update);s.mux.HandleFunc("DELETE /api/terms/{id}",s.del)
s.mux.HandleFunc("GET /api/search",s.search);s.mux.HandleFunc("GET /api/categories",s.categories);s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"terms":oe(s.db.List())})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){var t store.Term;json.NewDecoder(r.Body).Decode(&t);if t.Name==""{we(w,400,"name required");return};s.db.Create(&t);wj(w,201,t)}
func(s *Server)get(w http.ResponseWriter,r *http.Request){t:=s.db.Get(r.PathValue("id"));if t==nil{we(w,404,"not found");return};wj(w,200,t)}
func(s *Server)update(w http.ResponseWriter,r *http.Request){id:=r.PathValue("id");ex:=s.db.Get(id);if ex==nil{we(w,404,"not found");return};var t store.Term;json.NewDecoder(r.Body).Decode(&t);if t.Name==""{t.Name=ex.Name};s.db.Update(id,&t);wj(w,200,s.db.Get(id))}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)search(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"terms":oe(s.db.Search(r.URL.Query().Get("q")))})}
func(s *Server)categories(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"categories":oe(s.db.Categories())})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]int{"terms":s.db.Count(),"categories":len(s.db.Categories())})}
func(s *Server)health(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"status":"ok","service":"glossary","terms":s.db.Count()})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
