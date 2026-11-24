const NAV = {
  guest: [],
  client: [
    ["services", "Мои услуги"], ["payments", "Платежи"],
    ["sessions", "Занятия"], ["standards", "Нормативы"],
  ],
  trainer: [
    ["sessions", "Расписание"], ["enrollments", "Явка"],
    ["standards", "Нормативы"], ["permits", "Допуски"],
  ],
  admin: [
    ["services", "Услуги"], ["payments", "Платежи"], ["sessions", "Занятия"],
    ["enrollments", "Явка"], ["standards", "Нормативы"], ["users", "Пользователи"],
    ["tariffs", "Тарифы"], ["lanes", "Дорожки"], ["groups", "Группы"], ["permits", "Допуски"], ["journal", "Журнал"],
  ],
};
const ADDABLE = {
  admin: ["services", "payments", "sessions", "users", "tariffs", "lanes", "groups", "standards", "enrollments", "permits"],
  trainer: ["standards"],
  client: [],
};

const WORK = { 1: "свободное", 2: "индивидуальное", 3: "групповое", 4: "спортивное" };
const TARIFF_FORMAT = {
  free_swim_one_time: "свободное, разовое",
  free_swim_subscription: "свободное, абонемент",
  individual_one_time: "индивидуальное, разовое",
  individual_subscription: "индивидуальное, абонемент",
  group_one_time: "групповое, разовое",
  group_subscription: "групповое, абонемент",
  sport_pass: "спорт-пасс",
};
const USER_STATUS = { active: "активен", blocked: "заблокирован" };
const SERVICE_STATUS = { active: "активна", frozen: "заморожена", closed: "закрыта" };
const LANE_STATUS = { available: "доступна", closed: "закрыта" };
const GROUP_STATUS = { recruiting: "набор", training: "тренировки", finished: "завершена" };
const PAY_STATUS = { pending: "ожидает", completed: "оплачен", cancelled: "отменён" };
const METHOD = { cash: "наличные", card: "карта", transfer: "перевод" };
const QUAL = { instructor: "инструктор", trainer: "тренер", senior_trainer: "старший тренер" };
const RANK = { none: "без разряда", "3": "III разряд", "2": "II разряд", "1": "I разряд", cms: "КМС", ms: "МС" };
const SPEC = {
  freestyle: "кроль", breaststroke: "брасс", butterfly: "баттерфляй",
  medley: "комплекс", kids: "дети", aqua_aerobics: "аквааэробика",
};
const STYLE = {
  freestyle: "кроль", breaststroke: "брасс", butterfly: "баттерфляй",
  backstroke: "на спине", medley: "комплекс",
};
const OUTCOME = { passed: "сдан", failed: "не сдан" };
const RESULT = { success: "успешно", error: "ошибка" };
const ROLE = { admin: "администратор", trainer: "тренер", client: "клиент" };
const ACTION = { login: "вход", create: "создание", update: "изменение", delete: "удаление" };
const ERRORS = {
  "invalid credentials": "Неверный логин или пароль",
  unauthorized: "Нужно войти снова",
  "user is blocked": "Учётная запись заблокирована",
  "internal server error": "Ошибка на сервере",
  "user not found": "Пользователь не найден",
  "user already exists": "Такой логин уже есть",
  "client not found": "Клиент не найден",
  "trainer not found": "Тренер не найден",
  "tariff not found": "Тариф не найден",
  "tariff already exists": "Тариф с таким названием уже есть",
  "group not found": "Группа не найдена",
  "lane not found": "Дорожка не найдена",
  "service not found": "Услуга не найдена",
  "payment not found": "Платёж не найден",
  "session not found": "Занятие не найдено",
  "enrollment not found": "Запись не найдена",
  "enrollment already exists": "Услуга уже записана на это занятие",
  "standard not found": "Норматив не найден",
  "work format not found": "Формат не найден",
  "journal entry not found": "Запись журнала не найдена",
  "permit not found": "Допуск не найден",
  "permit already exists": "Такой допуск уже есть",
  "passport already exists": "Такой паспорт уже есть",
  forbidden: "Недостаточно прав",
  "service must be active": "Услуга не активна",
  "session date is outside the service period": "Дата занятия вне срока услуги",
  "visit limit is reached": "Лимит посещений исчерпан",
  "medical certificate is required": "Без медицинской справки услугу оформить нельзя",
  "service belongs to another client": "Это не ваша услуга",
  "client card is still in use": "Нельзя сменить роль: у клиента есть связанные записи",
  "trainer card is still in use": "Нельзя сменить роль: у тренера есть связанные записи",
  "passport is required": "Введите паспорт",
  "birth_date is required": "Укажите дату рождения",
  "birth_date must be in the past": "Дата рождения должна быть в прошлом",
  "qualification must be instructor, trainer or senior_trainer": "Некорректная квалификация",
  "rank must be none, 3, 2, 1, cms or ms": "Некорректный разряд",
  "experience_years must be at least 0": "Стаж не может быть отрицательным",
  "specialization must be freestyle, breaststroke, butterfly, medley, kids or aqua_aerobics": "Некорректная специализация",
  "work format must be individual, group or sport": "Допуск только к индивидуальному, групповому или спортивному формату",
  "issued_on is required": "Укажите дату выдачи",
  "expires_on must be on or after issued_on": "Срок допуска раньше даты выдачи",
  "login and password are required": "Введите логин и пароль",
  "login is required": "Введите логин",
  "login must not contain spaces": "В логине не должно быть пробелов",
  "login must be at most 50 characters": "Логин не длиннее 50 символов",
  "password must be at least 8 characters": "Пароль не короче 8 символов",
  "last_name is required": "Введите фамилию",
  "last_name must be at most 50 characters": "Фамилия не длиннее 50 символов",
  "first_name is required": "Введите имя",
  "first_name must be at most 50 characters": "Имя не длиннее 50 символов",
  "email is invalid": "Некорректная почта",
  "email must be at most 254 characters": "Почта не длиннее 254 символов",
  "role must be admin, trainer or client": "Некорректная роль",
  "status must be active or blocked": "Некорректный статус пользователя",
  "name is required": "Введите название",
  "name must be at most 100 characters": "Название не длиннее 100 символов",
  "name must be at most 20 characters": "Название не длиннее 20 символов",
  "start_date is required": "Укажите дату начала",
  "end_date is required": "Укажите дату окончания",
  "end_date must be on or after start_date": "Дата окончания раньше даты начала",
  "capacity must be greater than 0 and at most 32767": "Вместимость должна быть больше 0",
  "status must be recruiting, training or finished": "Некорректный статус группы",
  "length_m must be 25 or 50": "Длина дорожки: 25 или 50 м",
  "depth_m must be greater than 0 and at most 9.99": "Глубина должна быть больше 0 и не больше 9.99",
  "status must be available or closed": "Некорректный статус дорожки",
  "unknown format": "Неизвестный формат тарифа",
  "price must be between 0 and 999999.99": "Цена должна быть от 0 до 999999.99",
  "visit_limit must be greater than 0": "Лимит визитов должен быть больше 0",
  "duration_days must be greater than 0": "Срок должен быть больше 0 дней",
  "sport_pass must have price 0 and no visit limit": "У спорт-пасса цена 0 и нет лимита визитов",
  "client_id is required": "Выберите клиента",
  "tariff_id is required": "Выберите тариф",
  "group_id is required for group and sport tariffs": "Для группы и спорта нужна группа",
  "group_id is only allowed for group and sport tariffs": "Группа только у группового и спортивного тарифа",
  "status must be active, frozen or closed": "Некорректный статус услуги",
  "service_id is required": "Выберите услугу",
  "service_id and session_id are required": "Выберите услугу и занятие",
  "session_id is required": "Выберите занятие",
  "amount must be greater than 0 and at most 999999.99": "Сумма должна быть больше 0",
  "purpose must be at most 100 characters": "Назначение не длиннее 100 символов",
  "method must be cash, card or transfer": "Некорректный способ оплаты",
  "status must be pending, completed or cancelled": "Некорректный статус платежа",
  "sport_pass service cannot have payments": "По спорт-пассу платёж не вносится",
  "paid_on is required": "Укажите дату платежа",
  "paid_on must be YYYY-MM-DD": "Дата платежа в формате ГГГГ-ММ-ДД",
  "service_id is invalid": "Услуга не выбрана",
  "work_format_id must be 1, 2, 3 or 4": "Некорректный формат занятия",
  "session_date is required": "Укажите дату занятия",
  "start_time is required": "Укажите время начала",
  "end_time is required": "Укажите время окончания",
  "end_time must be after start_time": "Окончание должно быть позже начала",
  "group_id is required for group and sport sessions": "Для группы и спорта нужна группа",
  "group_id is only allowed for group and sport sessions": "Группа только у группового и спортивного занятия",
  "free swim session cannot have a trainer": "У свободного плавания не может быть тренера",
  "trainer_id is invalid": "Тренер не выбран",
  "trainer_id is required": "Выберите тренера",
  "topic must be at most 150 characters": "Тема не длиннее 150 символов",
  "service and session work formats differ": "Формат услуги и занятия не совпадает",
  "group service must be enrolled on a session of the same group": "Групповую услугу можно записать только на занятие той же группы",
  "comment must be at most 300 characters": "Комментарий не длиннее 300 символов",
  "style must be freestyle, breaststroke, butterfly, backstroke or medley": "Некорректный стиль",
  "distance_m must be 25, 50, 100, 200 or 400": "Дистанция: 25, 50, 100, 200 или 400 м",
  "held_on is required": "Укажите дату норматива",
  "outcome must be passed or failed": "Итог: сдан или не сдан",
  "attempt_no must be at least 1": "Номер попытки должен быть не меньше 1",
  "result_time must be HH:MM:SS": "Результат в формате ЧЧ:ММ:СС",
  "to must be on or after from": "Конец периода раньше начала",
  "to must be after from": "Конец периода не позже начала",
  "id is required": "Некорректный запрос",
  "filter id is invalid": "Некорректный фильтр",
};

let me = null;
const PAGE_SIZE = 20;

let page = "services";
let tablePage = 1;
let q = "";
let form = null;
let people = "accounts";
let filters = { action: "", from: "", to: "" };
let store = {};
let saving = false;
let renderToken = 0;
let pendingMsg = null;
let focusSearch = false;
let searchTimer = 0;
let bookID = 0;

const esc = (s) => String(s ?? "").replace(/[&<>"']/g, (ch) => ({
  "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;",
}[ch]));
const FIELD = {
  patronymic: "Отчество", phone: "Телефон", login: "Логин", last_name: "Фамилия",
  first_name: "Имя", email: "Почта", name: "Название", topic: "Тема",
  purpose: "Назначение", comment: "Комментарий", object: "Объект",
  description: "Описание", ip_address: "IP-адрес", passport: "Паспорт", address: "Адрес",
};
function explain(message, status) {
  const detail = message && message.includes(": ") ? message.slice(message.lastIndexOf(": ") + 2) : message;
  const known = (message && ERRORS[message]) || (detail && ERRORS[detail]);
  if (known) return known;
  if (detail && /^trainer \d+ has no permit/.test(detail)) return "У тренера нет допуска к этому формату на дату занятия";
  const atMost = detail && detail.match(/^([a-z_]+) must be at most (\d+) characters$/);
  if (atMost) return `${FIELD[atMost[1]] || "Поле"} не длиннее ${atMost[2]} символов`;
  if (detail && / must be YYYY-MM-DD$/.test(detail)) return "Дата в формате ГГГГ-ММ-ДД";
  if (detail && / must be HH:MM$/.test(detail)) return "Время в формате ЧЧ:ММ";
  if (detail && / is required$/.test(detail)) return "Заполните обязательное поле";
  if (detail && / is invalid$/.test(detail)) return "Некорректное значение";
  if (status === 400) return "Проверьте введённые данные";
  if (status === 401) return "Нужно войти снова";
  if (status === 403) return "Учётная запись заблокирована";
  if (status === 404) return "Не найдено";
  if (status === 409) return "Такая запись уже есть";
  if (status >= 500) return "Ошибка на сервере";
  return "Не удалось выполнить запрос";
}
const dash = (value) => value == null || value === "" ? "—" : value;
const tag = (text, warn) => `<span class="tag${warn ? " warn" : ""}">${esc(text)}</span>`;
const num = (value) => {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : 0;
};
const today = () => {
  const date = new Date();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${date.getFullYear()}-${month}-${day}`;
};
const money = (value) => Number(value).toFixed(2);
const hit = (parts) => !q.trim() || parts.join(" ").toLowerCase().includes(q.trim().toLowerCase());

function personName(user) {
  if (!user) return "—";
  const initials = [user.first_name, user.patronymic].filter(Boolean).map((part) => `${part[0]}.`).join("");
  return `${user.last_name} ${initials}`.trim();
}
function byID(rows) {
  return Object.fromEntries(rows.map((row) => [row.id, row]));
}
function clientName(data, clientID) {
  const client = data.clientByID[clientID];
  return client?.name || personName(client && data.userByID[client.user_id]);
}
function trainerName(data, trainerID) {
  const trainer = data.trainerByID[trainerID];
  return personName(trainer && data.userByID[trainer.user_id]);
}
function groupName(data, groupID) {
  return groupID && data.groupByID[groupID] ? data.groupByID[groupID].name : "—";
}
function laneName(data, laneID) {
  return laneID && data.laneByID[laneID] ? data.laneByID[laneID].name : "—";
}
function tariffName(data, tariffID) {
  const tariff = data.tariffByID[tariffID];
  return tariff ? tariff.name : "—";
}
function textName(value) {
  return value || "—";
}
function myClient() {
  return me && me.client_id ? { id: me.client_id, user_id: me.id } : null;
}
function myTrainer() {
  return me && me.trainer_id ? { id: me.trainer_id, user_id: me.id } : null;
}

function tbl(heads, rowHtml) {
  const rows = Array.isArray(rowHtml) ? rowHtml : [];
  const pages = Math.max(1, Math.ceil(rows.length / PAGE_SIZE));
  if (tablePage > pages) tablePage = pages;
  if (tablePage < 1) tablePage = 1;
  const start = (tablePage - 1) * PAGE_SIZE;
  const html = rows.slice(start, start + PAGE_SIZE).join("")
    || `<tr><td class="empty" colspan="${heads.length}">Ничего не найдено</td></tr>`;
  return `<div class="table-wrap"><table><thead><tr>${heads.map((head) => `<th>${head}</th>`).join("")}</tr></thead><tbody>${html}</tbody></table></div>${pager(rows.length, pages)}`;
}
function pager(total, pages) {
  if (total <= PAGE_SIZE) return "";
  const from = (tablePage - 1) * PAGE_SIZE + 1;
  const to = Math.min(total, tablePage * PAGE_SIZE);
  return `<div class="pager">
    <button type="button" class="btn-line" data-page="prev" ${tablePage === 1 ? "disabled" : ""}>назад</button>
    <span>${from}–${to} из ${total}</span>
    <button type="button" class="btn-line" data-page="next" ${tablePage === pages ? "disabled" : ""}>дальше</button>
  </div>`;
}
function formBtns() {
  return `<div class="wide"><button type="button" class="btn" data-act="save">Сохранить</button> <button type="button" class="btn-line" data-form="off">Отмена</button></div>`;
}
function searchBar(placeholder) {
  const add = (ADDABLE[me.role] || []).includes(page) && !(page === "users" && people !== "accounts");
  return `<div class="toolbar">
    <input id="q" placeholder="${esc(placeholder)}" value="${esc(q)}" />
    ${add ? `<button type="button" class="btn" data-form="new">Добавить</button>` : ""}
  </div>`;
}
function options(items, value, label, blank) {
  const head = blank == null ? "" : `<option value="">${esc(blank)}</option>`;
  return head + items.map((item) => `<option value="${esc(value(item))}">${esc(label(item))}</option>`).join("");
}
function select(name, items, value, label, blank) {
  return `<select name="${name}">${options(items, value, label, blank)}</select>`;
}

let refreshing = null;
async function api(path, opts = {}, retry = true) {
  const response = await fetch("/api/v1" + path, {
    method: opts.method || "GET",
    credentials: "same-origin",
    headers: opts.body ? { "Content-Type": "application/json" } : undefined,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
  });
  if (response.status === 401 && retry && path !== "/login" && path !== "/refresh") {
    refreshing ||= fetch("/api/v1/refresh", { method: "POST", credentials: "same-origin" })
      .finally(() => { refreshing = null; });
    if ((await refreshing).ok) return api(path, opts, false);
    me = null;
    store = {};
    const error = new Error("unauthorized");
    error.status = 401;
    throw error;
  }
  if (response.status === 204) return null;
  const text = await response.text();
  let data = null;
  if (text) {
    try { data = JSON.parse(text); } catch { data = { error: text }; }
  }
  if (!response.ok) {
    let message = data?.error || "";
    if (path === "/login" && response.status === 401) message = "invalid credentials";
    const error = new Error(message);
    error.status = response.status;
    throw error;
  }
  return data;
}
function load(key, path) {
  if (!store[key]) {
    store[key] = api(path).catch((error) => {
      delete store[key];
      throw error;
    });
  }
  return store[key];
}
function drop(prefix) {
  Object.keys(store).forEach((key) => {
    if (key === prefix || key.startsWith(prefix + ":")) delete store[key];
  });
}
function emptyRefs(clients = []) {
  return {
    users: [], clients, trainers: [], tariffs: [], groups: [], lanes: [],
    userByID: {}, clientByID: byID(clients), trainerByID: {},
    tariffByID: {}, groupByID: {}, laneByID: {},
  };
}
async function refs() {
  if (!me || me.role !== "admin") {
    const clients = me && me.role === "trainer" ? await load("clients", "/clients/") : [];
    return emptyRefs(clients);
  }
  const [users, clients, trainers, tariffs, groups, lanes] = await Promise.all([
    load("users", "/users/"),
    load("clients", "/clients/"),
    load("trainers", "/trainers/"),
    load("tariffs", "/tariffs/"),
    load("groups", "/groups/"),
    load("lanes", "/lanes/"),
  ]);
  return {
    users, clients, trainers, tariffs, groups, lanes,
    userByID: byID(users),
    clientByID: byID(clients),
    trainerByID: byID(trainers),
    tariffByID: byID(tariffs),
    groupByID: byID(groups),
    laneByID: byID(lanes),
  };
}

function loginPage() {
  return `<div class="login-page">
    <h1>Вход</h1>
    <div class="field"><label>Логин</label><input id="login" autocomplete="username" /></div>
    <div class="field"><label>Пароль</label><input id="password" type="password" autocomplete="current-password" /></div>
    <button type="button" class="btn" data-enter>Войти</button>
  </div>`;
}

async function servicesPage() {
  const data = await refs();
  let rows = await load("services", "/services/");
  rows = rows.filter((row) => hit([
    row.id, row.client_name, row.tariff_name, row.group_name,
    SERVICE_STATUS[row.status], row.start_date, row.end_date,
  ]));
  const heads = ["№", "Клиент", "Тариф", "Группа", "Статус", "Срок"].concat(me.role === "admin" ? [""] : []);
  return `<h2>${me.role === "client" ? "Мои услуги" : "Услуги"}</h2>
    ${searchBar("Поиск: клиент, тариф, группа")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Клиент</label>${select("client_id", data.clients, (item) => item.id, (item) => item.name || clientName(data, item.id))}</div>
      <div><label>Тариф</label>${select("tariff_id", data.tariffs, (item) => item.id, (item) => item.name)}</div>
      <div><label>Группа</label>${select("group_id", data.groups, (item) => item.id, (item) => item.name, "—")}</div>
      <div><label>Начало</label><input name="start_date" type="date" value="${today()}" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(heads, rows.map((row) => `<tr>
      <td>${row.id}</td><td class="name">${esc(textName(row.client_name))}</td>
      <td>${esc(textName(row.tariff_name))}</td><td>${esc(textName(row.group_name))}</td>
      <td>${tag(SERVICE_STATUS[row.status] || row.status, row.status === "frozen")}</td>
      <td>${esc(row.start_date)} — ${esc(row.end_date)}</td>
      ${me.role === "admin" ? `<td class="row-actions">${serviceActions(row)}</td>` : ""}
    </tr>`))}`;
}
function serviceActions(row) {
  if (row.status === "closed") return "";
  const links = [];
  if (row.status === "active") links.push(statusLink(`/services/${row.id}/status`, "frozen", "services", "заморозить"));
  if (row.status === "frozen") links.push(statusLink(`/services/${row.id}/status`, "active", "services", "активировать"));
  links.push(statusLink(`/services/${row.id}/status`, "closed", "services", "закрыть"));
  return links.join("");
}

async function paymentsPage() {
  const data = await refs();
  const services = me.role === "admin" ? await load("services:all", "/services/") : [];
  let rows = await load("payments", "/payments/");
  const payable = services.filter((item) => data.tariffByID[item.tariff_id]?.format !== "sport_pass");
  rows = rows.filter((row) => hit([
    row.client_name, `услуга ${row.service_id}`, money(row.amount),
    row.paid_on, METHOD[row.method], row.purpose, PAY_STATUS[row.status],
  ]));
  return `<h2>Платежи</h2>
    ${searchBar("Поиск: клиент, сумма, назначение")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Услуга</label>${select("service_id", payable, (item) => item.id, (item) => `${item.id} — ${item.client_name || clientName(data, item.client_id)}`)}</div>
      <div><label>Сумма</label><input name="amount" value="12.00" /></div>
      <div><label>Способ</label><select name="method">${Object.entries(METHOD).map(([code, label]) => `<option value="${code}"${code === "card" ? " selected" : ""}>${label}</option>`).join("")}</select></div>
      <div><label>Статус</label><select name="status">${Object.entries(PAY_STATUS).map(([code, label]) => `<option value="${code}"${code === "completed" ? " selected" : ""}>${label}</option>`).join("")}</select></div>
      <div><label>Назначение</label><input name="purpose" value="Оплата тарифа" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Клиент", "Услуга", "Сумма", "Дата", "Способ", "Назначение", "Статус"].concat(me.role === "admin" ? [""] : []),
      rows.map((row) => {
        return `<tr>
          <td class="name">${esc(textName(row.client_name))}</td>
          <td>услуга ${row.service_id}</td><td class="money">${money(row.amount)}</td>
          <td>${esc(row.paid_on)}</td><td>${esc(METHOD[row.method] || row.method)}</td>
          <td>${esc(dash(row.purpose))}</td>
          <td>${tag(PAY_STATUS[row.status] || row.status, row.status !== "completed")}</td>
          ${me.role === "admin" ? `<td class="row-actions">${paymentActions(row)}</td>` : ""}
        </tr>`;
      }))}`;
}
function paymentActions(row) {
  if (row.status !== "pending") return "";
  return statusLink(`/payments/${row.id}/status`, "completed", "payments", "провести")
    + statusLink(`/payments/${row.id}/status`, "cancelled", "payments", "отменить");
}

async function sessionsPage() {
  const data = await refs();
  let rows = await load("sessions", "/sessions/");
  const trainer = myTrainer();
  if (me.role === "trainer") rows = rows.filter((row) => trainer && row.trainer_ids.includes(trainer.id));
  let ownServices = [];
  if (me.role === "client") {
    ownServices = await load("services:mine", "/services/");
    const groups = new Set(ownServices.map((item) => item.group_id).filter(Boolean));
    rows = rows.filter((row) => row.work_format_id === 1 || (row.group_id && groups.has(row.group_id)));
  }
  const listed = rows.slice();
  rows = rows.filter((row) => hit([
    row.session_date, row.start_time, row.end_time, WORK[row.work_format_id],
    row.group_name, row.lane_name, row.topic, (row.trainer_names || []).join(" "),
    row.was_held ? "проведено" : "запланировано",
  ]));
  const booking = bookID ? bookForm(listed, ownServices) : "";
  return `<h2>${me.role === "trainer" ? "Расписание" : "Занятия"}</h2>
    ${searchBar("Поиск: дата, формат, дорожка")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Дата</label><input name="session_date" type="date" value="${today()}" /></div>
      <div><label>Начало</label><input name="start_time" value="09:00" /></div>
      <div><label>Конец</label><input name="end_time" value="09:45" /></div>
      <div><label>Формат</label><select name="work_format_id">${Object.entries(WORK).map(([code, label]) => `<option value="${code}">${label}</option>`).join("")}</select></div>
      <div><label>Дорожка</label>${select("lane_id", data.lanes, (item) => item.id, (item) => item.name, "—")}</div>
      <div><label>Группа</label>${select("group_id", data.groups, (item) => item.id, (item) => item.name, "—")}</div>
      <div><label>Тема</label><input name="topic" /></div>
      <div><label>Тренер</label>${select("trainer_id", data.trainers, (item) => item.id, (item) => trainerName(data, item.id), "—")}</div>
      ${formBtns()}
    </div>` : ""}
    ${booking}
    ${tbl(["Дата", "Время", "Формат", "Группа", "Дорожка", "Тема", "Тренер", "Статус"].concat(me.role === "client" || me.role === "admin" ? [""] : []),
      rows.map((row) => `<tr>
        <td class="name">${esc(row.session_date)}</td>
        <td>${esc(row.start_time)}–${esc(row.end_time)}</td>
        <td>${esc(WORK[row.work_format_id] || row.work_format_id)}</td>
        <td>${esc(textName(row.group_name))}</td>
        <td>${esc(textName(row.lane_name))}</td>
        <td>${esc(dash(row.topic))}</td>
        <td>${esc((row.trainer_names || []).length ? row.trainer_names.join(", ") : "—")}</td>
        <td>${tag(row.was_held ? "проведено" : "запланировано")}</td>
        ${me.role === "admin" ? `<td class="row-actions"><a href="#" data-act="held" data-path="/sessions/${row.id}/held" data-held="${row.was_held ? "0" : "1"}">${row.was_held ? "отменить" : "провести"}</a></td>` : ""}
        ${me.role === "client" ? `<td class="row-actions"><a href="#" data-act="book" data-id="${row.id}">записаться</a></td>` : ""}
      </tr>`))}`;
}
function bookForm(rows, services) {
  const session = rows.find((row) => row.id === bookID);
  if (!session) return "";
  const matches = services.filter((item) => {
    if (item.status !== "active" || item.work_format_id !== session.work_format_id) return false;
    return session.work_format_id < 3 || item.group_id === session.group_id;
  });
  if (!matches.length) {
    return `<p class="msg">Нет подходящей услуги для записи</p><button type="button" class="btn-line" data-form="off">Отмена</button>`;
  }
  return `<div class="form-grid">
    <div class="wide"><label>Запись на ${esc(session.session_date)} ${esc(session.start_time)}</label>
      <select name="service_id">${matches.map((item) => `<option value="${item.id}">${esc(item.tariff_name || "услуга " + item.id)}</option>`).join("")}</select>
    </div>
    <div class="wide"><label>Комментарий</label><input name="comment" /></div>
    ${formBtns()}
  </div>`;
}

async function enrollmentsPage() {
  const data = await refs();
  const [items, sessions] = await Promise.all([
    load("enrollments", "/enrollments/"),
    load("sessions", "/sessions/"),
  ]);
  const services = me.role === "admin" ? await load("services:all", "/services/") : [];
  const sessionByID = byID(sessions);
  const trainer = myTrainer();
  let rows = items.filter((row) => {
    const session = sessionByID[row.session_id];
    if (me.role === "trainer" && (!trainer || !session || !session.trainer_ids.includes(trainer.id))) return false;
    return hit([
      sessionLabel(session), row.client_name, row.tariff_name,
      row.attended ? "был" : "не был", row.comment,
    ]);
  });
  return `<h2>Явка</h2>
    ${searchBar("Поиск: клиент, занятие")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Занятие</label><select name="session_id">${options(sessions, (item) => item.id, (item) => sessionLabel(item))}</select></div>
      <div><label>Услуга</label><select name="service_id">${services.map((item) => {
        return `<option value="${item.id}" data-format="${item.work_format_id || ""}" data-group="${item.group_id || ""}">${esc(item.id + " — " + (item.client_name || clientName(data, item.client_id)))}</option>`;
      }).join("")}</select></div>
      <div class="wide"><label>Комментарий</label><input name="comment" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Занятие", "Клиент", "Тариф", "Явка", "Комментарий"].concat(me.role !== "client" ? [""] : []),
      rows.map((row) => {
        const session = sessionByID[row.session_id];
        return `<tr>
          <td class="name">${esc(sessionLabel(session))}</td>
          <td>${esc(textName(row.client_name))}</td>
          <td>${esc(textName(row.tariff_name))}</td>
          <td>${tag(row.attended ? "был" : "не был", !row.attended)}</td>
          <td>${esc(dash(row.comment))}</td>
          ${me.role !== "client" ? `<td class="row-actions"><a href="#" data-act="attend" data-path="/enrollments/${row.service_id}/${row.session_id}" data-attended="${row.attended ? "0" : "1"}">${row.attended ? "снять явку" : "отметить явку"}</a></td>` : ""}
        </tr>`;
      }))}`;
}
function sessionLabel(session) {
  if (!session) return "—";
  return `${session.session_date} ${session.start_time}`;
}

async function standardsPage() {
  const data = await refs();
  const trainer = myTrainer();
  let rows = await load("standards", "/standards/");
  rows = rows.filter((row) => hit([
    row.client_name, row.trainer_name,
    STYLE[row.style], row.distance_m, row.held_on, row.result_time, OUTCOME[row.outcome], row.attempt_no,
  ]));
  const trainerField = me.role === "trainer"
    ? `<input type="hidden" name="trainer_id" value="${trainer ? trainer.id : ""}" />`
    : `<div><label>Тренер</label>${select("trainer_id", data.trainers, (item) => item.id, (item) => trainerName(data, item.id))}</div>`;
  return `<h2>Нормативы</h2>
    ${searchBar("Поиск: клиент, стиль")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Клиент</label>${select("client_id", data.clients, (item) => item.id, (item) => item.name || clientName(data, item.id))}</div>
      ${trainerField}
      <div><label>Стиль</label><select name="style">${Object.entries(STYLE).map(([code, label]) => `<option value="${code}">${label}</option>`).join("")}</select></div>
      <div><label>Дистанция, м</label><select name="distance_m"><option>25</option><option>50</option><option>100</option><option>200</option><option>400</option></select></div>
      <div><label>Дата</label><input name="held_on" type="date" value="${today()}" /></div>
      <div><label>Результат</label><input name="result_time" placeholder="00:00:26" /></div>
      <div><label>Итог</label><select name="outcome">${Object.entries(OUTCOME).map(([code, label]) => `<option value="${code}">${label}</option>`).join("")}</select></div>
      <div><label>Попытка</label><input name="attempt_no" value="1" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Клиент", "Тренер", "Стиль", "м", "Дата", "Время", "Итог", "Попытка"],
      rows.map((row) => `<tr>
        <td class="name">${esc(textName(row.client_name))}</td>
        <td>${esc(textName(row.trainer_name))}</td>
        <td>${esc(STYLE[row.style] || row.style)}</td><td>${row.distance_m}</td>
        <td>${esc(row.held_on)}</td><td class="money">${esc(dash(row.result_time))}</td>
        <td>${tag(OUTCOME[row.outcome] || row.outcome, row.outcome === "failed")}</td>
        <td>${row.attempt_no}</td>
      </tr>`))}`;
}

async function usersPage() {
  const data = await refs();
  const tabs = [["accounts", "Учётные записи"], ["clients", "Клиенты"], ["trainers", "Тренеры"]];
  const sub = `<div class="subnav">${tabs.map(([id, title]) => `<a href="#" data-s="${id}" class="${people === id ? "on" : ""}">${title}</a>`).join("")}</div>`;
  if (people === "clients") {
    const rows = data.clients.filter((row) => {
      const user = data.userByID[row.user_id];
      return hit([personName(user), row.passport, row.birth_date, row.address, row.medical_certificate ? "есть" : "нет"]);
    });
    return `<h2>Пользователи</h2>
      ${sub}${searchBar("Поиск: ФИО, паспорт")}
      ${tbl(["ФИО", "Паспорт", "Рождение", "Адрес", "Медсправка", ""], rows.map((row) => `<tr>
        <td class="name">${esc(row.name || personName(data.userByID[row.user_id]))}</td>
        <td>${esc(row.passport)}</td><td>${esc(row.birth_date)}</td><td>${esc(dash(row.address))}</td>
        <td>${tag(row.medical_certificate ? "есть" : "нет", !row.medical_certificate)}</td>
        <td class="row-actions"><a href="#" data-act="medical" data-id="${row.id}" data-on="${row.medical_certificate ? "0" : "1"}">${row.medical_certificate ? "отозвать" : "справка есть"}</a></td>
      </tr>`))}`;
  }
  if (people === "trainers") {
    const permits = await load("permits", "/permits/");
    const formats = {};
    permits.forEach((item) => {
      (formats[item.trainer_id] ||= []).push(WORK[item.work_format_id] || item.work_format_id);
    });
    const rows = data.trainers.filter((row) => hit([
      trainerName(data, row.id), QUAL[row.qualification], RANK[row.rank], row.experience_years, SPEC[row.specialization],
    ]));
    return `<h2>Пользователи</h2>
      ${sub}${searchBar("Поиск: ФИО, квалификация")}
      ${tbl(["ФИО", "Квалификация", "Разряд", "Стаж", "Специализация", "Допуски"], rows.map((row) => `<tr>
        <td class="name">${esc(trainerName(data, row.id))}</td>
        <td>${esc(QUAL[row.qualification] || row.qualification)}</td>
        <td>${esc(RANK[row.rank] || row.rank)}</td>
        <td>${row.experience_years}</td>
        <td>${esc(SPEC[row.specialization] || row.specialization)}</td>
        <td>${esc(formats[row.id]?.join(", ") || "—")}</td>
      </tr>`))}`;
  }
  const rows = data.users.filter((row) => hit([row.login, personName(row), row.role, row.phone, row.email, USER_STATUS[row.status]]));
  return `<h2>Пользователи</h2>
    ${sub}${searchBar("Поиск: логин, фамилия")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Логин</label><input name="login" /></div>
      <div><label>Пароль</label><input name="password" type="password" placeholder="не короче 8 символов" /></div>
      <div><label>Роль</label><select name="role">${Object.entries(ROLE).map(([code, label]) => `<option value="${code}"${code === "client" ? " selected" : ""}>${label}</option>`).join("")}</select></div>
      <div><label>Фамилия</label><input name="last_name" /></div>
      <div><label>Имя</label><input name="first_name" /></div>
      <div><label>Отчество</label><input name="patronymic" /></div>
      <div><label>Телефон</label><input name="phone" /></div>
      <div><label>Почта</label><input name="email" /></div>
      ${profileFields()}
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Логин", "ФИО", "Роль", "Телефон", "Почта", "Статус", ""], rows.map((row) => `<tr>
      <td>${esc(row.login)}</td><td class="name">${esc(personName(row))}</td><td>${esc(ROLE[row.role] || row.role)}</td>
      <td>${esc(dash(row.phone))}</td><td>${esc(dash(row.email))}</td>
      <td>${tag(USER_STATUS[row.status] || row.status, row.status === "blocked")}</td>
      <td class="row-actions">${statusLink(`/users/${row.id}/status`, row.status === "active" ? "blocked" : "active", "users", row.status === "active" ? "заблокировать" : "разблокировать")}</td>
    </tr>`))}`;
}

async function tariffsPage() {
  let rows = await load("tariffs", "/tariffs/");
  rows = rows.filter((row) => hit([row.id, row.name, TARIFF_FORMAT[row.format], row.price, row.visit_limit, row.duration_days]));
  return `<h2>Тарифы</h2>
    ${searchBar("Поиск: название, формат")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Название</label><input name="name" /></div>
      <div><label>Формат</label><select name="format">${Object.entries(TARIFF_FORMAT).map(([code, label]) => `<option value="${code}">${label}</option>`).join("")}</select></div>
      <div><label>Цена</label><input name="price" value="0.00" /></div>
      <div><label>Визитов</label><input name="visit_limit" /></div>
      <div><label>Дней</label><input name="duration_days" value="30" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["№", "Название", "Формат", "Цена", "Визиты", "Дни"], rows.map((row) => `<tr>
      <td>${row.id}</td><td class="name">${esc(row.name)}</td><td>${esc(TARIFF_FORMAT[row.format] || row.format)}</td>
      <td class="money">${money(row.price)}</td><td>${dash(row.visit_limit)}</td><td>${row.duration_days}</td>
    </tr>`))}`;
}

async function lanesPage() {
  let rows = await load("lanes", "/lanes/");
  rows = rows.filter((row) => hit([row.id, row.name, row.length_m, row.depth_m, LANE_STATUS[row.status]]));
  return `<h2>Дорожки</h2>
    ${searchBar("Поиск: название, статус")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Название</label><input name="name" value="Дорожка 9" /></div>
      <div><label>Длина, м</label><select name="length_m"><option value="25">25</option><option value="50">50</option></select></div>
      <div><label>Глубина, м</label><input name="depth_m" value="1.80" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["№", "Название", "Длина", "Глубина", "Статус", ""], rows.map((row) => `<tr>
      <td>${row.id}</td><td class="name">${esc(row.name)}</td><td>${row.length_m} м</td>
      <td>${Number(row.depth_m).toFixed(2)} м</td>
      <td>${tag(LANE_STATUS[row.status] || row.status, row.status === "closed")}</td>
      <td class="row-actions">${statusLink(`/lanes/${row.id}/status`, row.status === "available" ? "closed" : "available", "lanes", row.status === "available" ? "закрыть" : "открыть")}</td>
    </tr>`))}`;
}

async function groupsPage() {
  const data = await refs();
  let rows = data.groups.filter((row) => hit([row.name, GROUP_STATUS[row.status], row.capacity, row.start_date, row.end_date]));
  return `<h2>Группы</h2>
    ${searchBar("Поиск: название, статус")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Название</label><input name="name" /></div>
      <div><label>Вместимость</label><input name="capacity" value="12" /></div>
      <div><label>Начало</label><input name="start_date" type="date" value="${today()}" /></div>
      <div><label>Конец</label><input name="end_date" type="date" value="${today()}" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Название", "Статус", "Места", "Срок"].concat(me.role === "admin" ? [""] : []), rows.map((row) => `<tr>
      <td class="name">${esc(row.name)}</td>
      <td>${tag(GROUP_STATUS[row.status] || row.status, row.status === "recruiting")}</td>
      <td class="money">${row.capacity}</td>
      <td>${esc(row.start_date)} — ${esc(row.end_date)}</td>
      ${me.role === "admin" ? `<td class="row-actions">${groupActions(row)}</td>` : ""}
    </tr>`))}`;
}
function groupActions(row) {
  if (row.status === "recruiting") return statusLink(`/groups/${row.id}/status`, "training", "groups", "начать тренировки");
  if (row.status === "training") return statusLink(`/groups/${row.id}/status`, "finished", "groups", "завершить");
  return "";
}

function profileFields() {
  const options = (dict) => Object.entries(dict).map(([code, label]) => `<option value="${code}">${label}</option>`).join("");
  return `
    <div data-for="client"><label>Паспорт</label><input name="passport" maxlength="20" /></div>
    <div data-for="client"><label>Дата рождения</label><input name="birth_date" type="date" /></div>
    <div data-for="client" class="wide"><label>Адрес</label><input name="address" /></div>
    <div data-for="client"><label>Медсправка</label><select name="medical_certificate"><option value="false">нет</option><option value="true">есть</option></select></div>
    <div data-for="trainer" hidden><label>Квалификация</label><select name="qualification">${options(QUAL)}</select></div>
    <div data-for="trainer" hidden><label>Разряд</label><select name="rank">${options(RANK)}</select></div>
    <div data-for="trainer" hidden><label>Стаж, лет</label><input name="experience_years" value="0" /></div>
    <div data-for="trainer" hidden><label>Специализация</label><select name="specialization">${options(SPEC)}</select></div>`;
}
function showRoleFields(role) {
  document.querySelectorAll("[data-for]").forEach((field) => {
    field.hidden = field.dataset.for !== role;
  });
}

async function permitsPage() {
  const data = await refs();
  let rows = await load("permits", "/permits/");
  rows = rows.filter((row) => hit([
    me.role === "trainer" ? personName(me) : trainerName(data, row.trainer_id),
    WORK[row.work_format_id], row.issued_on, row.expires_on,
  ]));
  const formats = [2, 3, 4].map((id) => `<option value="${id}">${WORK[id]}</option>`).join("");
  return `<h2>Допуски</h2>
    ${searchBar("Поиск: тренер, формат")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Тренер</label>${select("trainer_id", data.trainers, (item) => item.id, (item) => trainerName(data, item.id))}</div>
      <div><label>Формат</label><select name="work_format_id">${formats}</select></div>
      <div><label>Выдан</label><input name="issued_on" type="date" value="${today()}" /></div>
      <div><label>До</label><input name="expires_on" type="date" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Тренер", "Формат", "Выдан", "До"], rows.map((row) => `<tr>
      <td class="name">${esc(me.role === "trainer" ? personName(me) : trainerName(data, row.trainer_id))}</td>
      <td>${esc(WORK[row.work_format_id] || row.work_format_id)}</td>
      <td>${esc(row.issued_on)}</td><td>${esc(row.expires_on || "бессрочно")}</td>
    </tr>`))}`;
}

async function journalPage() {
  const params = new URLSearchParams();
  if (filters.action) params.set("action", filters.action);
  if (filters.from) params.set("from", filters.from);
  if (filters.to) params.set("to", filters.to);
  if (q.trim()) params.set("q", q.trim());
  const rows = await load("journal:" + params.toString(), "/journal/?" + params.toString());
  return `<h2>Журнал</h2>
    <div class="filters">
      <input id="q" placeholder="поиск по описанию" value="${esc(q)}" />
      <select data-filter="action">
        <option value="">все действия</option>
        ${Object.entries(ACTION).map(([code, label]) => `<option value="${code}"${filters.action === code ? " selected" : ""}>${label}</option>`).join("")}
      </select>
      <input data-filter="from" type="date" value="${esc(filters.from)}" />
      <input data-filter="to" type="date" value="${esc(filters.to)}" />
    </div>
    ${tbl(["Время", "Кто", "Действие", "Объект", "Описание", "IP", "Результат"], rows.map((row) => `<tr>
      <td>${esc(row.occurred_at)}</td><td class="name">${esc(row.login)}</td><td>${esc(ACTION[row.action] || row.action)}</td>
      <td>${esc(row.object)}</td><td>${esc(row.description)}</td><td>${esc(row.ip_address)}</td>
      <td>${tag(RESULT[row.result] || row.result, row.result === "error")}</td>
    </tr>`))}`;
}

function statusLink(path, status, bucket, label) {
  return `<a href="#" data-act="status" data-path="${path}" data-status="${status}" data-drop="${bucket}">${label}</a>`;
}
function themeControl() {
  const dark = document.body.classList.contains("theme-dark");
  return `<button type="button" class="theme" data-theme aria-pressed="${dark}" aria-label="Тема">
    <span>тема</span>
    <b>
      <svg class="sun" viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="4"/><path d="M12 2.5v3M12 18.5v3M2.5 12h3M18.5 12h3M5.2 5.2l2.1 2.1M16.7 16.7l2.1 2.1M5.2 18.8l2.1-2.1M16.7 7.3l2.1-2.1"/></svg>
      <svg class="moon" viewBox="0 0 24 24" aria-hidden="true"><path d="M16 3.5A7.5 7.5 0 1 0 20.5 15 5.8 5.8 0 0 1 16 3.5z"/></svg>
    </b>
  </button>`;
}
function shell(main) {
  const links = me ? NAV[me.role].map(([id, title]) => `<a href="#" data-p="${id}" class="${page === id ? "on" : ""}">${title}</a>`).join("") : "";
  const account = me ? `<div class="account"><span class="account-name">${esc(personName(me))}</span><span class="account-role">${esc(ROLE[me.role] || me.role)}</span><button type="button" class="icon-btn" data-out aria-label="Выйти" title="Выйти"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M10 6H6v12h4"/><path d="M10 12h9"/><path d="M15.5 8.5 19 12l-3.5 3.5"/></svg></button></div>` : "";
  return `<header class="top">${themeControl()}${account}</header><aside><span class="brand">Бассейн</span><nav>${links}</nav></aside><main>${main}</main>`;
}
function showMsg(message, status) {
  let el = document.querySelector("main .msg");
  if (!el) {
    el = document.createElement("p");
    el.className = "msg";
    const head = document.querySelector("main h1, main h2");
    if (head) head.insertAdjacentElement("afterend", el);
    else document.querySelector("main")?.prepend(el);
  }
  el.textContent = explain(message, status);
}
function sessionEnded(error) {
  return error.status === 401 || error.message === "user is blocked";
}
function report(error) {
  const inside = !!document.querySelector("aside nav a");
  if (sessionEnded(error) && inside) {
    me = null;
    store = {};
    pendingMsg = error;
    render();
    return;
  }
  showMsg(error.message, error.status);
}

const pages = {
  services: servicesPage, payments: paymentsPage, sessions: sessionsPage,
  enrollments: enrollmentsPage, standards: standardsPage, users: usersPage,
  tariffs: tariffsPage, lanes: lanesPage, groups: groupsPage,
  permits: permitsPage, journal: journalPage,
};

async function render() {
  const token = ++renderToken;
  const app = document.getElementById("app");
  app.className = "layout" + (me ? "" : " login-full");
  if (!me) {
    app.innerHTML = shell(loginPage());
    if (pendingMsg) {
      const notice = pendingMsg;
      pendingMsg = null;
      showMsg(notice.message, notice.status);
    }
    return;
  }
  app.innerHTML = shell(`<p class="empty">Загрузка…</p>`);
  let html;
  try {
    html = await pages[page]();
  } catch (error) {
    if (sessionEnded(error)) {
      me = null;
      store = {};
      pendingMsg = error;
      if (token === renderToken) render();
      return;
    }
    html = `<h2>Ошибка</h2><p class="msg">${esc(explain(error.message, error.status))}</p>`;
  }
  if (token !== renderToken) return;
  app.innerHTML = shell(html);
  if (focusSearch) {
    const field = document.getElementById("q");
    if (field) {
      field.focus();
      field.setSelectionRange(q.length, q.length);
    }
    focusSearch = false;
  }
  if (page === "enrollments" && form === "new") filterEnrollmentServices();
}

async function filterEnrollmentServices() {
  const sessionField = document.querySelector("[name=session_id]");
  const serviceField = document.querySelector("[name=service_id]");
  if (!sessionField || !serviceField) return;
  const sessions = await load("sessions", "/sessions/");
  const session = sessions.find((item) => item.id === num(sessionField.value));
  let first = "";
  [...serviceField.options].forEach((option) => {
    const sameFormat = session && Number(option.dataset.format) === session.work_format_id;
    const sameGroup = session && (session.work_format_id < 3 || option.dataset.group === String(session.group_id || ""));
    option.hidden = !(sameFormat && sameGroup);
    if (!option.hidden && !first) first = option.value;
  });
  if (serviceField.selectedOptions[0]?.hidden) serviceField.value = first;
}

function fields() {
  const data = {};
  document.querySelectorAll(".form-grid [name]").forEach((field) => {
    data[field.name] = field.value;
  });
  return data;
}
async function save() {
  if (saving) return;
  saving = true;
  const input = fields();
  try {
    if (page === "services") {
      const tariffs = await load("tariffs", "/tariffs/");
      const tariff = tariffs.find((item) => item.id === num(input.tariff_id));
      const needsGroup = tariff && (tariff.work_format_id === 3 || tariff.work_format_id === 4);
      await api("/services/", { method: "POST", body: {
        client_id: num(input.client_id),
        tariff_id: num(input.tariff_id),
        group_id: needsGroup ? num(input.group_id) : null,
        start_date: input.start_date,
      }});
      drop("services");
    } else if (page === "payments") {
      await api("/payments/", { method: "POST", body: {
        service_id: num(input.service_id),
        amount: Number(input.amount),
        method: input.method,
        purpose: input.purpose,
        status: input.status,
      }});
      drop("payments");
    } else if (page === "sessions" && bookID) {
      await api("/enrollments/", { method: "POST", body: {
        service_id: num(input.service_id),
        session_id: bookID,
        comment: input.comment || "",
      }});
      bookID = 0;
      drop("enrollments");
    } else if (page === "sessions") {
      const format = num(input.work_format_id);
      const needsGroup = format === 3 || format === 4;
      await api("/sessions/", { method: "POST", body: {
        work_format_id: format,
        group_id: needsGroup ? num(input.group_id) : null,
        lane_id: input.lane_id ? num(input.lane_id) : null,
        session_date: input.session_date,
        start_time: input.start_time,
        end_time: input.end_time,
        topic: input.topic || null,
        trainer_ids: format === 1 || !input.trainer_id ? [] : [num(input.trainer_id)],
      }});
      drop("sessions");
    } else if (page === "enrollments") {
      await api("/enrollments/", { method: "POST", body: {
        service_id: num(input.service_id),
        session_id: num(input.session_id),
        comment: input.comment || "",
      }});
      drop("enrollments");
    } else if (page === "standards") {
      await api("/standards/", { method: "POST", body: {
        client_id: num(input.client_id),
        trainer_id: num(input.trainer_id),
        style: input.style,
        distance_m: num(input.distance_m),
        held_on: input.held_on,
        result_time: input.result_time || "",
        outcome: input.outcome,
        attempt_no: num(input.attempt_no) || 1,
      }});
      drop("standards");
    } else if (page === "permits") {
      await api("/permits/", { method: "POST", body: {
        trainer_id: num(input.trainer_id),
        work_format_id: num(input.work_format_id),
        issued_on: input.issued_on,
        expires_on: input.expires_on || "",
      }});
      drop("permits");
    } else if (page === "users") {
      await api("/users/", { method: "POST", body: {
        login: input.login,
        password: input.password,
        role: input.role,
        last_name: input.last_name,
        first_name: input.first_name,
        patronymic: input.patronymic || "",
        phone: input.phone || "",
        email: input.email || "",
        passport: input.passport || "",
        birth_date: input.birth_date || "",
        address: input.address || "",
        medical_certificate: input.medical_certificate === "true",
        qualification: input.qualification || "",
        rank: input.rank || "",
        experience_years: num(input.experience_years),
        specialization: input.specialization || "",
      }});
      drop("users");
      drop("clients");
      drop("trainers");
    } else if (page === "tariffs") {
      const sport = input.format === "sport_pass";
      await api("/tariffs/", { method: "POST", body: {
        name: input.name,
        format: input.format,
        price: sport ? 0 : Number(input.price),
        visit_limit: sport || input.visit_limit === "" ? null : num(input.visit_limit),
        duration_days: num(input.duration_days),
      }});
      drop("tariffs");
    } else if (page === "lanes") {
      await api("/lanes/", { method: "POST", body: {
        name: input.name,
        length_m: num(input.length_m),
        depth_m: Number(input.depth_m),
        status: "available",
      }});
      drop("lanes");
    } else if (page === "groups") {
      await api("/groups/", { method: "POST", body: {
        name: input.name,
        capacity: num(input.capacity),
        start_date: input.start_date,
        end_date: input.end_date,
      }});
      drop("groups");
    }
    form = null;
    render();
  } catch (error) {
    report(error);
  } finally {
    saving = false;
  }
}

document.getElementById("app").addEventListener("click", async (event) => {
  const theme = event.target.closest("[data-theme]");
  if (theme) {
    const dark = !document.body.classList.contains("theme-dark");
    document.body.classList.toggle("theme-dark", dark);
    document.body.classList.toggle("theme-light", !dark);
    theme.setAttribute("aria-pressed", String(dark));
    return;
  }
  const step = event.target.closest("[data-page]");
  if (step) {
    event.preventDefault();
    if (step.disabled) return;
    tablePage += step.dataset.page === "next" ? 1 : -1;
    render();
    return;
  }
  const nav = event.target.closest("a[data-p]");
  if (nav) {
    event.preventDefault();
    page = nav.dataset.p;
    form = null;
    bookID = 0;
    q = "";
    people = "accounts";
    tablePage = 1;
    render();
    return;
  }
  const sub = event.target.closest("a[data-s]");
  if (sub) {
    event.preventDefault();
    people = sub.dataset.s;
    form = null;
    q = "";
    tablePage = 1;
    render();
    return;
  }
  if (event.target.closest("[data-out]")) {
    try { await api("/logout", { method: "POST" }); } catch { /* куки всё равно сбрасываем локально */ }
    me = null;
    store = {};
    form = null;
    q = "";
    render();
    return;
  }
  if (event.target.closest("[data-enter]")) {
    try {
      me = await api("/login", { method: "POST", body: {
        login: document.getElementById("login").value,
        password: document.getElementById("password").value,
      }});
      store = {};
      form = null;
      q = "";
      page = NAV[me.role][0][0];
      render();
    } catch (error) {
      report(error);
    }
    return;
  }
  if (event.target.closest("[data-form=new]")) { form = "new"; render(); return; }
  if (event.target.closest("[data-form=off]")) { form = null; bookID = 0; render(); return; }
  if (event.target.closest("[data-act=save]")) { event.preventDefault(); save(); return; }

  const act = event.target.closest("[data-act]");
  if (!act) return;
  event.preventDefault();
  try {
    if (act.dataset.act === "status") {
      await api(act.dataset.path, { method: "PATCH", body: { status: act.dataset.status } });
      drop(act.dataset.drop);
    } else if (act.dataset.act === "held") {
      await api(act.dataset.path, { method: "PATCH", body: { was_held: act.dataset.held === "1" } });
      drop("sessions");
    } else if (act.dataset.act === "attend") {
      await api(act.dataset.path, { method: "PATCH", body: { attended: act.dataset.attended === "1" } });
      drop("enrollments");
    } else if (act.dataset.act === "book") {
      bookID = num(act.dataset.id);
      form = null;
    } else if (act.dataset.act === "medical") {
      await api(`/clients/${act.dataset.id}/medical`, {
        method: "PATCH",
        body: { medical_certificate: act.dataset.on === "1" },
      });
      drop("clients");
    }
    render();
  } catch (error) {
    report(error);
  }
});
document.getElementById("app").addEventListener("input", (event) => {
  if (event.target.id !== "q") return;
  q = event.target.value;
  tablePage = 1;
  focusSearch = true;
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => render(), page === "journal" ? 250 : 0);
});
document.getElementById("app").addEventListener("change", (event) => {
  if (event.target.name === "session_id") filterEnrollmentServices();
  if (event.target.name === "role") showRoleFields(event.target.value);
  if (!event.target.dataset.filter) return;
  filters[event.target.dataset.filter] = event.target.value;
  tablePage = 1;
  drop("journal");
  render();
});

api("/me").then((user) => {
  me = user;
  page = NAV[user.role][0][0];
  render();
}).catch(() => render());
