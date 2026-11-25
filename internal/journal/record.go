package journal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"pool/internal/auth"
)

func Record(svc *Service) func(*gin.Context) {
	return func(c *gin.Context) {
		record(svc, c)
	}
}

func record(svc *Service, c *gin.Context) {
	status := c.Writer.Status()
	path := c.Request.URL.Path
	method := c.Request.Method
	if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
		return
	}
	if path == "/api/v1/refresh" || path == "/api/v1/logout" {
		return
	}
	u, ok := auth.CurrentUser(c)
	if !ok {
		return
	}
	action := ActionUpdate
	switch method {
	case http.MethodPost:
		action = ActionCreate
	case http.MethodDelete:
		action = ActionDelete
	}
	object := resource(path)
	isLogin := path == "/api/v1/login"
	if isLogin {
		action = ActionLogin
		object = "users"
	}
	result := ResultSuccess
	if status >= http.StatusBadRequest {
		result = ResultError
	}
	_, _ = svc.Create(c.Request.Context(), Entry{
		UserID:      u.ID,
		Action:      action,
		Object:      object,
		Description: describe(c, action, object, path, u.Login, status),
		Result:      result,
		IPAddress:   c.ClientIP(),
	})
}

func resource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 && parts[2] != "" {
		return parts[2]
	}
	return "api"
}

func describe(c *gin.Context, action, object, path, login string, status int) string {
	if action == ActionLogin {
		who := login
		if who == "" {
			who = "?"
		}
		if status >= http.StatusBadRequest {
			if status == http.StatusForbidden {
				return clipDesc("вход отклонён: учётка заблокирована (" + who + ")")
			}
			return clipDesc("вход отклонён: неверные данные (" + who + ")")
		}
		return clipDesc("вход в систему (" + who + ")")
	}
	head := actionLabel(action) + " " + objectLabel(object)
	ids := pathIDs(path, object)
	if ids != "" {
		head += " " + ids
	}
	if extra := pathExtra(path); extra != "" {
		head += ", " + extra
	}
	fields := bodySummary(auth.RequestBody(c))
	if fields != "" {
		head += ": " + fields
	}
	if status >= http.StatusBadRequest {
		head += " — " + failReason(status)
	}
	return clipDesc(head)
}

func failReason(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "ошибка данных"
	case http.StatusUnauthorized:
		return "нет авторизации"
	case http.StatusForbidden:
		return "запрещено"
	case http.StatusNotFound:
		return "не найдено"
	case http.StatusConflict:
		return "конфликт"
	case http.StatusInternalServerError:
		return "ошибка сервера"
	default:
		return fmt.Sprintf("отказ (%d)", status)
	}
}

func actionLabel(action string) string {
	switch action {
	case ActionCreate:
		return "создание"
	case ActionUpdate:
		return "изменение"
	case ActionDelete:
		return "удаление"
	default:
		return action
	}
}

func objectLabel(object string) string {
	switch object {
	case "sessions":
		return "занятия"
	case "services":
		return "услуги"
	case "users":
		return "пользователя"
	case "payments":
		return "платежа"
	case "enrollments":
		return "явки"
	case "standards":
		return "норматива"
	case "permits":
		return "допуска"
	case "tariffs":
		return "тарифа"
	case "lanes":
		return "дорожки"
	case "groups":
		return "группы"
	case "clients":
		return "клиента"
	case "trainers":
		return "тренера"
	default:
		return object
	}
}

func pathIDs(path, object string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 4 {
		return ""
	}
	rest := parts[3:]
	switch object {
	case "enrollments":
		if len(rest) >= 2 {
			return "услуга " + rest[0] + ", занятие " + rest[1]
		}
	case "permits":
		if len(rest) >= 2 {
			return "тренер " + rest[0] + ", формат " + formatName(rest[1])
		}
		if len(rest) >= 1 {
			return "тренер " + rest[0]
		}
	}
	if looksID(rest[0]) {
		return rest[0]
	}
	return ""
}

func pathExtra(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 5 {
		return ""
	}
	last := parts[len(parts)-1]
	if looksID(last) {
		return ""
	}
	switch last {
	case "status":
		return "статус"
	case "password":
		return "пароль"
	case "held":
		return "проведение"
	case "medical":
		return "медсправка"
	default:
		return last
	}
}

func looksID(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func bodySummary(raw []byte) string {
	raw = bytesTrim(raw)
	if len(raw) == 0 || raw[0] != '{' {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}
	var parts []string
	for _, key := range fieldOrder {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		if skipField[key] {
			continue
		}
		s := formatValue(key, v)
		if s == "" {
			continue
		}
		parts = append(parts, fieldLabel(key)+" "+s)
	}
	return strings.Join(parts, ", ")
}

var skipField = map[string]bool{
	"password": true,
}

var fieldOrder = []string{
	"login", "last_name", "first_name", "patronymic", "role", "status", "phone", "email",
	"session_date", "start_time", "end_time", "work_format_id", "was_held",
	"lane_id", "group_id", "client_id", "trainer_id", "trainer_ids",
	"name", "format", "capacity", "price", "duration_months",
	"service_id", "amount", "method", "paid_at",
	"attended", "held_on", "style", "distance_m", "outcome",
	"issued_on", "expires_on", "medical_certificate",
	"number",
}

func fieldLabel(key string) string {
	switch key {
	case "login":
		return "логин"
	case "last_name":
		return "фамилия"
	case "first_name":
		return "имя"
	case "patronymic":
		return "отчество"
	case "role":
		return "роль"
	case "status":
		return "статус"
	case "phone":
		return "телефон"
	case "email":
		return "почта"
	case "session_date":
		return "дата"
	case "start_time":
		return "начало"
	case "end_time":
		return "конец"
	case "work_format_id":
		return "формат"
	case "was_held":
		return "статус занятия"
	case "lane_id":
		return "дорожка"
	case "group_id":
		return "группа"
	case "client_id":
		return "клиент"
	case "trainer_id":
		return "тренер"
	case "trainer_ids":
		return "тренеры"
	case "name":
		return "название"
	case "format":
		return "формат"
	case "capacity":
		return "ёмкость"
	case "price":
		return "цена"
	case "duration_months":
		return "срок, мес"
	case "service_id":
		return "услуга"
	case "amount":
		return "сумма"
	case "method":
		return "способ"
	case "paid_at":
		return "оплачено"
	case "attended":
		return "явка"
	case "held_on":
		return "дата"
	case "style":
		return "стиль"
	case "distance_m":
		return "дистанция, м"
	case "outcome":
		return "результат"
	case "issued_on":
		return "выдан"
	case "expires_on":
		return "до"
	case "medical_certificate":
		return "медсправка"
	case "number":
		return "номер"
	default:
		return key
	}
}

func formatValue(key string, v any) string {
	switch key {
	case "work_format_id":
		return formatName(fmt.Sprint(asInt(v)))
	case "was_held":
		if asBool(v) {
			return "проведено"
		}
		return "запланировано"
	case "attended":
		if asBool(v) {
			return "присутствовал"
		}
		return "отсутствовал"
	case "medical_certificate":
		if asBool(v) {
			return "есть"
		}
		return "нет"
	case "trainer_ids":
		return joinAny(v)
	}
	switch t := v.(type) {
	case string:
		s := strings.TrimSpace(t)
		if i := strings.IndexByte(s, 'T'); i > 0 && looksDatePrefix(s) {
			return s[:i]
		}
		if len(s) >= 8 && s[2] == ':' {
			return s
		}
		return s
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "да"
		}
		return "нет"
	case json.Number:
		return t.String()
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(b)
	}
}

func formatName(id string) string {
	switch id {
	case "1":
		return "свободное"
	case "2":
		return "индивидуальное"
	case "3":
		return "групповое"
	case "4":
		return "спортивное"
	default:
		return id
	}
}

func asInt(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case json.Number:
		n, _ := t.Int64()
		return int(n)
	case int:
		return t
	case string:
		n, _ := strconv.Atoi(t)
		return n
	default:
		n, _ := strconv.Atoi(fmt.Sprint(v))
		return n
	}
}

func asBool(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return t == "true" || t == "1"
	case float64:
		return t != 0
	default:
		return false
	}
}

func joinAny(v any) string {
	switch t := v.(type) {
	case []any:
		parts := make([]string, 0, len(t))
		for _, item := range t {
			parts = append(parts, fmt.Sprint(asInt(item)))
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprint(v)
	}
}

func looksDatePrefix(s string) bool {
	return len(s) >= 10 && s[4] == '-' && s[7] == '-'
}

func bytesTrim(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}

func clipDesc(s string) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= 300 {
		return s
	}
	runes := []rune(s)
	return string(runes[:300])
}
