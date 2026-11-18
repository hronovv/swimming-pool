const NAV = {
  guest: [],
  client: [
    ["services", "Мои услуги"], ["payments", "Платежи"],
    ["sessions", "Занятия"], ["standards", "Нормативы"],
  ],
  trainer: [
    ["sessions", "Расписание"], ["enrollments", "Явка"],
    ["groups", "Группы"], ["standards", "Нормативы"], ["permits", "Допуски"],
  ],
  admin: [
    ["services", "Услуги"], ["payments", "Платежи"], ["sessions", "Занятия"],
    ["enrollments", "Явка"], ["standards", "Нормативы"], ["users", "Пользователи"],
    ["tariffs", "Тарифы"], ["lanes", "Дорожки"], ["groups", "Группы"], ["journal", "Журнал"],
  ],
};
const ADDABLE = {
  admin: ["services", "payments", "sessions", "users", "tariffs", "lanes", "groups", "standards", "enrollments"],
  trainer: ["standards"],
  client: [],
};

const WORK = { 1: "свободное", 2: "индивид", 3: "группа", 4: "спорт" };
const TARIFF_FORMAT = {
  free_swim_one_time: "свободное разовое",
  free_swim_subscription: "свободное абонемент",
  individual_one_time: "индивид разовое",
  individual_subscription: "индивид абонемент",
  group_one_time: "группа разовое",
  group_subscription: "группа абонемент",
  sport_pass: "спорт-пасс",
};
const USER_STATUS = { active: "активен", blocked: "заблокирован" };
const SERVICE_STATUS = { active: "активна", frozen: "заморожена", closed: "закрыта" };
const LANE_STATUS = { available: "доступна", closed: "закрыта" };
const GROUP_STATUS = { recruiting: "набор", training: "тренировки", finished: "завершена" };
const PAY_STATUS = { pending: "ожидает", completed: "выполнен", cancelled: "отменён" };
const METHOD = { cash: "наличные", card: "карта", transfer: "перевод" };
const QUAL = { instructor: "инструктор", trainer: "тренер", senior_trainer: "старший тренер" };
const RANK = { none: "нет", "3": "3 разряд", "2": "2 разряд", "1": "1 разряд", cms: "кмс", ms: "мс" };
const SPEC = {
  freestyle: "кроль", breaststroke: "брасс", butterfly: "баттерфляй",
  medley: "комплекс", kids: "дети", aqua_aerobics: "аквааэробика",
};
const STYLE = {
  freestyle: "кроль", breaststroke: "брасс", butterfly: "баттерфляй",
  backstroke: "на спине", medley: "комплекс",
};
const OUTCOME = { passed: "сдан", failed: "не сдан" };
const RESULT = { success: "успех", error: "ошибка" };
const ROLE = { admin: "админ", trainer: "тренер", client: "клиент" };
const ERRORS = {
  "invalid credentials": "Неверный логин или пароль",
  unauthorized: "Нужно войти снова",
  "user is blocked": "Учётная запись заблокирована",
  "enrollment already exists": "Услуга уже записана на это занятие",
};

let me = null;
let page = "services";
let q = "";
let form = null;
let people = "accounts";
let filters = { action: "", from: "", to: "" };
let store = {};
let saving = false;
let renderToken = 0;
let focusSearch = false;
let searchTimer = 0;

const esc = (s) => String(s ?? "").replace(/[&<>"']/g, (ch) => ({
  "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;",
}[ch]));
const explain = (message) => ERRORS[message] || message || "Ошибка запроса";
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
  return personName(client && data.userByID[client.user_id]);
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
function myClient(data) {
  return data.clients.find((item) => item.user_id === me.id) || null;
}
function myTrainer(data) {
  return data.trainers.find((item) => item.user_id === me.id) || null;
}

function tbl(heads, body) {
  const html = body || `<tr><td class="empty" colspan="${heads.length}">Ничего не найдено</td></tr>`;
  return `<div class="table-wrap"><table><thead><tr>${heads.map((head) => `<th>${head}</th>`).join("")}</tr></thead><tbody>${html}</tbody></table></div>`;
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
    const error = new Error(data?.error || "ошибка запроса");
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
async function refs() {
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
    <p class="lead">Логин и пароль учётной записи бассейна.</p>
    <div class="field"><label>Логин</label><input id="login" autocomplete="username" /></div>
    <div class="field"><label>Пароль</label><input id="password" type="password" autocomplete="current-password" /></div>
    <button type="button" class="btn" data-enter>Войти</button>
  </div>`;
}

async function servicesPage() {
  const data = await refs();
  const mine = myClient(data);
  const path = me.role === "client" && mine ? `/services/?client_id=${mine.id}` : "/services/";
  let rows = me.role === "client" && !mine ? [] : await load("services:" + path, path);
  rows = rows.filter((row) => hit([
    row.id, clientName(data, row.client_id), tariffName(data, row.tariff_id),
    groupName(data, row.group_id), SERVICE_STATUS[row.status], row.start_date, row.end_date,
  ]));
  const heads = ["id", "Клиент", "Тариф", "Группа", "Статус", "Срок"].concat(me.role === "admin" ? [""] : []);
  return `<h2>${me.role === "client" ? "Мои услуги" : "Услуги"}</h2>
    <p class="lead">${me.role === "client" ? "Оформленные тарифы." : "Оформление услуги клиенту, смена статуса."}</p>
    ${searchBar("Поиск: клиент, тариф, группа")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Клиент</label>${select("client_id", data.clients, (item) => item.id, (item) => clientName(data, item.id))}</div>
      <div><label>Тариф</label>${select("tariff_id", data.tariffs, (item) => item.id, (item) => item.name)}</div>
      <div><label>Группа</label>${select("group_id", data.groups, (item) => item.id, (item) => item.name, "—")}</div>
      <div><label>Начало</label><input name="start_date" type="date" value="${today()}" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(heads, rows.map((row) => `<tr>
      <td>${row.id}</td><td class="name">${esc(clientName(data, row.client_id))}</td>
      <td>${esc(tariffName(data, row.tariff_id))}</td><td>${esc(groupName(data, row.group_id))}</td>
      <td>${tag(SERVICE_STATUS[row.status] || row.status, row.status === "frozen")}</td>
      <td>${esc(row.start_date)} — ${esc(row.end_date)}</td>
      ${me.role === "admin" ? `<td class="row-actions">${serviceActions(row)}</td>` : ""}
    </tr>`).join(""))}`;
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
  const mine = myClient(data);
  let services = [];
  if (me.role !== "client") services = await load("services:all", "/services/");
  else if (mine) services = await load("services:mine", `/services/?client_id=${mine.id}`);
  const serviceByID = byID(services);
  let rows = await load("payments", "/payments/");
  if (me.role === "client") {
    const ids = new Set(services.map((item) => item.id));
    rows = rows.filter((row) => ids.has(row.service_id));
  }
  const payable = services.filter((item) => data.tariffByID[item.tariff_id]?.format !== "sport_pass");
  rows = rows.filter((row) => {
    const service = serviceByID[row.service_id];
    return hit([
      clientName(data, service?.client_id), `услуга ${row.service_id}`, money(row.amount),
      row.paid_on, METHOD[row.method], row.purpose, PAY_STATUS[row.status],
    ]);
  });
  return `<h2>Платежи</h2>
    <p class="lead">${me.role === "client" ? "Оплаты по своим услугам." : "Внесение оплаты по услуге. Спорт-пасс без платежа."}</p>
    ${searchBar("Поиск: клиент, сумма, назначение")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Услуга</label>${select("service_id", payable, (item) => item.id, (item) => `${item.id} — ${clientName(data, item.client_id)}`)}</div>
      <div><label>Сумма</label><input name="amount" value="12.00" /></div>
      <div><label>Способ</label><select name="method"><option value="card">карта</option><option value="cash">наличные</option><option value="transfer">перевод</option></select></div>
      <div><label>Статус</label><select name="status"><option value="completed">выполнен</option><option value="pending">ожидает</option><option value="cancelled">отменён</option></select></div>
      <div><label>Назначение</label><input name="purpose" value="Оплата тарифа" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Клиент", "Услуга", "Сумма", "Дата", "Способ", "Назначение", "Статус"].concat(me.role === "admin" ? [""] : []),
      rows.map((row) => {
        const service = serviceByID[row.service_id];
        return `<tr>
          <td class="name">${esc(clientName(data, service?.client_id))}</td>
          <td>услуга ${row.service_id}</td><td class="money">${money(row.amount)}</td>
          <td>${esc(row.paid_on)}</td><td>${esc(METHOD[row.method] || row.method)}</td>
          <td>${esc(dash(row.purpose))}</td>
          <td>${tag(PAY_STATUS[row.status] || row.status, row.status !== "completed")}</td>
          ${me.role === "admin" ? `<td class="row-actions">${paymentActions(row)}</td>` : ""}
        </tr>`;
      }).join(""))}`;
}
function paymentActions(row) {
  if (row.status !== "pending") return "";
  return statusLink(`/payments/${row.id}/status`, "completed", "payments", "выполнить")
    + statusLink(`/payments/${row.id}/status`, "cancelled", "payments", "отменить");
}

async function sessionsPage() {
  const data = await refs();
  let rows = await load("sessions", "/sessions/");
  const trainer = myTrainer(data);
  if (me.role === "trainer") rows = rows.filter((row) => trainer && row.trainer_ids.includes(trainer.id));
  if (me.role === "client") {
    const mine = myClient(data);
    const services = mine ? await load("services:mine", `/services/?client_id=${mine.id}`) : [];
    const groups = new Set(services.map((item) => item.group_id).filter(Boolean));
    rows = rows.filter((row) => row.work_format_id === 1 || (row.group_id && groups.has(row.group_id)));
  }
  rows = rows.filter((row) => hit([
    row.session_date, row.start_time, row.end_time, WORK[row.work_format_id],
    groupName(data, row.group_id), laneName(data, row.lane_id), row.topic,
    row.trainer_ids.map((id) => trainerName(data, id)).join(" "),
    row.was_held ? "проведено" : "запланировано",
  ]));
  return `<h2>${me.role === "trainer" ? "Расписание" : "Занятия"}</h2>
    <p class="lead">${me.role === "trainer" ? "Свои занятия." : me.role === "client" ? "Ближайшие занятия." : "Сессии: формат, группа, дорожка, тренер."}</p>
    ${searchBar("Поиск: дата, формат, дорожка")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Дата</label><input name="session_date" type="date" value="${today()}" /></div>
      <div><label>Начало</label><input name="start_time" value="09:00" /></div>
      <div><label>Конец</label><input name="end_time" value="09:45" /></div>
      <div><label>Формат</label><select name="work_format_id"><option value="1">свободное</option><option value="2">индивид</option><option value="3">группа</option><option value="4">спорт</option></select></div>
      <div><label>Дорожка</label>${select("lane_id", data.lanes, (item) => item.id, (item) => item.name, "—")}</div>
      <div><label>Группа</label>${select("group_id", data.groups, (item) => item.id, (item) => item.name, "—")}</div>
      <div><label>Тема</label><input name="topic" /></div>
      <div><label>Тренер</label>${select("trainer_id", data.trainers, (item) => item.id, (item) => trainerName(data, item.id), "—")}</div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Дата", "Время", "Формат", "Группа", "Дорожка", "Тема", "Тренер", "Статус"].concat(me.role === "admin" ? [""] : []),
      rows.map((row) => `<tr>
        <td class="name">${esc(row.session_date)}</td>
        <td>${esc(row.start_time)}–${esc(row.end_time)}</td>
        <td>${esc(WORK[row.work_format_id] || row.work_format_id)}</td>
        <td>${esc(groupName(data, row.group_id))}</td>
        <td>${esc(laneName(data, row.lane_id))}</td>
        <td>${esc(dash(row.topic))}</td>
        <td>${esc(row.trainer_ids.length ? row.trainer_ids.map((id) => trainerName(data, id)).join(", ") : "—")}</td>
        <td>${tag(row.was_held ? "проведено" : "запланировано")}</td>
        ${me.role === "admin" ? `<td class="row-actions"><a href="#" data-act="held" data-path="/sessions/${row.id}/held" data-held="${row.was_held ? "0" : "1"}">${row.was_held ? "снять" : "проведено"}</a></td>` : ""}
      </tr>`).join(""))}`;
}

async function enrollmentsPage() {
  const data = await refs();
  const [items, sessions, services] = await Promise.all([
    load("enrollments", "/enrollments/"),
    load("sessions", "/sessions/"),
    load("services:all", "/services/"),
  ]);
  const sessionByID = byID(sessions);
  const serviceByID = byID(services);
  const trainer = myTrainer(data);
  let rows = items.filter((row) => {
    const session = sessionByID[row.session_id];
    const service = serviceByID[row.service_id];
    if (me.role === "trainer" && (!trainer || !session || !session.trainer_ids.includes(trainer.id))) return false;
    return hit([
      sessionLabel(session), clientName(data, service?.client_id), tariffName(data, service?.tariff_id),
      row.attended ? "да" : "нет", row.comment,
    ]);
  });
  return `<h2>Явка</h2>
    <p class="lead">Запись услуги на занятие и отметка посещения.</p>
    ${searchBar("Поиск: клиент, занятие")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Занятие</label><select name="session_id">${options(sessions, (item) => item.id, (item) => sessionLabel(item))}</select></div>
      <div><label>Услуга</label><select name="service_id">${services.map((item) => {
        const tariff = data.tariffByID[item.tariff_id];
        return `<option value="${item.id}" data-format="${tariff?.work_format_id || ""}" data-group="${item.group_id || ""}">${esc(item.id + " — " + clientName(data, item.client_id))}</option>`;
      }).join("")}</select></div>
      <div class="wide"><label>Комментарий</label><input name="comment" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Занятие", "Клиент", "Тариф", "Явка", "Комментарий"].concat(me.role !== "client" ? [""] : []),
      rows.map((row) => {
        const session = sessionByID[row.session_id];
        const service = serviceByID[row.service_id];
        return `<tr>
          <td class="name">${esc(sessionLabel(session))}</td>
          <td>${esc(clientName(data, service?.client_id))}</td>
          <td>${esc(tariffName(data, service?.tariff_id))}</td>
          <td>${tag(row.attended ? "да" : "нет", !row.attended)}</td>
          <td>${esc(dash(row.comment))}</td>
          ${me.role !== "client" ? `<td class="row-actions"><a href="#" data-act="attend" data-path="/enrollments/${row.service_id}/${row.session_id}" data-attended="${row.attended ? "0" : "1"}">${row.attended ? "снять явку" : "отметить"}</a></td>` : ""}
        </tr>`;
      }).join(""))}`;
}
function sessionLabel(session) {
  if (!session) return "—";
  return `${session.session_date} ${session.start_time}`;
}

async function standardsPage() {
  const data = await refs();
  const mine = myClient(data);
  const trainer = myTrainer(data);
  const path = me.role === "client" && mine ? `/standards/?client_id=${mine.id}` : "/standards/";
  let rows = me.role === "client" && !mine ? [] : await load("standards:" + path, path);
  rows = rows.filter((row) => hit([
    clientName(data, row.client_id), trainerName(data, row.trainer_id),
    STYLE[row.style], row.distance_m, row.held_on, row.result_time, OUTCOME[row.outcome], row.attempt_no,
  ]));
  const trainerField = me.role === "trainer"
    ? `<input type="hidden" name="trainer_id" value="${trainer ? trainer.id : ""}" />`
    : `<div><label>Тренер</label>${select("trainer_id", data.trainers, (item) => item.id, (item) => trainerName(data, item.id))}</div>`;
  return `<h2>Нормативы</h2>
    <p class="lead">${me.role === "client" ? "Свои попытки." : "Стиль, дистанция, время, сдан / не сдан."}</p>
    ${searchBar("Поиск: клиент, стиль")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Клиент</label>${select("client_id", data.clients, (item) => item.id, (item) => clientName(data, item.id))}</div>
      ${trainerField}
      <div><label>Стиль</label><select name="style">${Object.entries(STYLE).map(([code, label]) => `<option value="${code}">${label}</option>`).join("")}</select></div>
      <div><label>Дистанция, м</label><select name="distance_m"><option>25</option><option>50</option><option>100</option><option>200</option><option>400</option></select></div>
      <div><label>Дата</label><input name="held_on" type="date" value="${today()}" /></div>
      <div><label>Результат</label><input name="result_time" placeholder="00:00:26" /></div>
      <div><label>Итог</label><select name="outcome"><option value="passed">сдан</option><option value="failed">не сдан</option></select></div>
      <div><label>Попытка</label><input name="attempt_no" value="1" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Клиент", "Тренер", "Стиль", "м", "Дата", "Время", "Итог", "Попытка"],
      rows.map((row) => `<tr>
        <td class="name">${esc(clientName(data, row.client_id))}</td>
        <td>${esc(trainerName(data, row.trainer_id))}</td>
        <td>${esc(STYLE[row.style] || row.style)}</td><td>${row.distance_m}</td>
        <td>${esc(row.held_on)}</td><td class="money">${esc(dash(row.result_time))}</td>
        <td>${tag(OUTCOME[row.outcome] || row.outcome, row.outcome === "failed")}</td>
        <td>${row.attempt_no}</td>
      </tr>`).join(""))}`;
}

async function usersPage() {
  const data = await refs();
  const tabs = [["accounts", "Счета"], ["clients", "Клиенты"], ["trainers", "Тренеры"]];
  const sub = `<div class="subnav">${tabs.map(([id, title]) => `<a href="#" data-s="${id}" class="${people === id ? "on" : ""}">${title}</a>`).join("")}</div>`;
  if (people === "clients") {
    const rows = data.clients.filter((row) => {
      const user = data.userByID[row.user_id];
      return hit([personName(user), row.passport, row.birth_date, row.address, row.medical_certificate ? "есть" : "нет"]);
    });
    return `<h2>Пользователи</h2>
      <p class="lead">Карточка клиента: паспорт, дата рождения, медсправка.</p>
      ${sub}${searchBar("Поиск: ФИО, паспорт")}
      ${tbl(["ФИО", "Паспорт", "Рождение", "Адрес", "Медсправка"], rows.map((row) => `<tr>
        <td class="name">${esc(personName(data.userByID[row.user_id]))}</td>
        <td>${esc(row.passport)}</td><td>${esc(row.birth_date)}</td><td>${esc(dash(row.address))}</td>
        <td>${tag(row.medical_certificate ? "есть" : "нет", !row.medical_certificate)}</td>
      </tr>`).join(""))}`;
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
      <p class="lead">Карточка тренера и допуски к форматам.</p>
      ${sub}${searchBar("Поиск: ФИО, квалификация")}
      ${tbl(["ФИО", "Квалификация", "Разряд", "Стаж", "Специализация", "Допуски"], rows.map((row) => `<tr>
        <td class="name">${esc(trainerName(data, row.id))}</td>
        <td>${esc(QUAL[row.qualification] || row.qualification)}</td>
        <td>${esc(RANK[row.rank] || row.rank)}</td>
        <td>${row.experience_years}</td>
        <td>${esc(SPEC[row.specialization] || row.specialization)}</td>
        <td>${esc(formats[row.id]?.join(", ") || "—")}</td>
      </tr>`).join(""))}`;
  }
  const rows = data.users.filter((row) => hit([row.login, personName(row), row.role, row.phone, row.email, USER_STATUS[row.status]]));
  return `<h2>Пользователи</h2>
    <p class="lead">Учётные записи: login, роль, блокировка.</p>
    ${sub}${searchBar("Поиск: логин, фамилия")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Логин</label><input name="login" /></div>
      <div><label>Пароль</label><input name="password" type="password" placeholder="не короче 8 символов" /></div>
      <div><label>Роль</label><select name="role"><option value="client">клиент</option><option value="trainer">тренер</option><option value="admin">админ</option></select></div>
      <div><label>Фамилия</label><input name="last_name" /></div>
      <div><label>Имя</label><input name="first_name" /></div>
      <div><label>Отчество</label><input name="patronymic" /></div>
      <div><label>Телефон</label><input name="phone" /></div>
      <div><label>Почта</label><input name="email" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Логин", "ФИО", "Роль", "Телефон", "Почта", "Статус", ""], rows.map((row) => `<tr>
      <td>${esc(row.login)}</td><td class="name">${esc(personName(row))}</td><td>${esc(ROLE[row.role] || row.role)}</td>
      <td>${esc(dash(row.phone))}</td><td>${esc(dash(row.email))}</td>
      <td>${tag(USER_STATUS[row.status] || row.status, row.status === "blocked")}</td>
      <td class="row-actions">${statusLink(`/users/${row.id}/status`, row.status === "active" ? "blocked" : "active", "users", row.status === "active" ? "заблокировать" : "активировать")}</td>
    </tr>`).join(""))}`;
}

async function tariffsPage() {
  let rows = await load("tariffs", "/tariffs/");
  rows = rows.filter((row) => hit([row.id, row.name, TARIFF_FORMAT[row.format], row.price, row.visit_limit, row.duration_days]));
  return `<h2>Тарифы</h2>
    <p class="lead">Справочник: формат, цена, лимит посещений, срок.</p>
    ${searchBar("Поиск: название, формат")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Название</label><input name="name" /></div>
      <div><label>Формат</label><select name="format">${Object.entries(TARIFF_FORMAT).map(([code, label]) => `<option value="${code}">${label}</option>`).join("")}</select></div>
      <div><label>Цена</label><input name="price" value="0.00" /></div>
      <div><label>Визитов</label><input name="visit_limit" /></div>
      <div><label>Дней</label><input name="duration_days" value="30" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["id", "Название", "Формат", "Цена", "Визитов", "Дней"], rows.map((row) => `<tr>
      <td>${row.id}</td><td class="name">${esc(row.name)}</td><td>${esc(TARIFF_FORMAT[row.format] || row.format)}</td>
      <td class="money">${money(row.price)}</td><td>${dash(row.visit_limit)}</td><td>${row.duration_days}</td>
    </tr>`).join(""))}`;
}

async function lanesPage() {
  let rows = await load("lanes", "/lanes/");
  rows = rows.filter((row) => hit([row.id, row.name, row.length_m, row.depth_m, LANE_STATUS[row.status]]));
  return `<h2>Дорожки</h2>
    <p class="lead">Длина 25 или 50 м, глубина, доступна / закрыта.</p>
    ${searchBar("Поиск: название, статус")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Название</label><input name="name" value="Дорожка 9" /></div>
      <div><label>Длина, м</label><select name="length_m"><option value="25">25</option><option value="50">50</option></select></div>
      <div><label>Глубина, м</label><input name="depth_m" value="1.80" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["id", "Название", "Длина", "Глубина", "Статус", ""], rows.map((row) => `<tr>
      <td>${row.id}</td><td class="name">${esc(row.name)}</td><td>${row.length_m} м</td>
      <td>${Number(row.depth_m).toFixed(2)} м</td>
      <td>${tag(LANE_STATUS[row.status] || row.status, row.status === "closed")}</td>
      <td class="row-actions">${statusLink(`/lanes/${row.id}/status`, row.status === "available" ? "closed" : "available", "lanes", row.status === "available" ? "закрыть" : "открыть")}</td>
    </tr>`).join(""))}`;
}

async function groupsPage() {
  const data = await refs();
  let rows = data.groups.filter((row) => hit([row.name, GROUP_STATUS[row.status], row.capacity, row.start_date, row.end_date]));
  return `<h2>Группы</h2>
    <p class="lead">Набор, тренировки, завершена. Вместимость и срок.</p>
    ${searchBar("Поиск: название, статус")}
    ${form === "new" && me.role === "admin" ? `<div class="form-grid">
      <div><label>Название</label><input name="name" /></div>
      <div><label>Вместимость</label><input name="capacity" value="12" /></div>
      <div><label>Начало</label><input name="start_date" type="date" value="${today()}" /></div>
      <div><label>Конец</label><input name="end_date" type="date" value="${today()}" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Название", "Статус", "Мест", "Срок"].concat(me.role === "admin" ? [""] : []), rows.map((row) => `<tr>
      <td class="name">${esc(row.name)}</td>
      <td>${tag(GROUP_STATUS[row.status] || row.status, row.status === "recruiting")}</td>
      <td class="money">${row.capacity}</td>
      <td>${esc(row.start_date)} — ${esc(row.end_date)}</td>
      ${me.role === "admin" ? `<td class="row-actions">${groupActions(row)}</td>` : ""}
    </tr>`).join(""))}`;
}
function groupActions(row) {
  if (row.status === "recruiting") return statusLink(`/groups/${row.id}/status`, "training", "groups", "к тренировкам");
  if (row.status === "training") return statusLink(`/groups/${row.id}/status`, "finished", "groups", "завершить");
  return "";
}

async function permitsPage() {
  const data = await refs();
  const trainer = myTrainer(data);
  let rows = await load("permits", "/permits/");
  if (me.role === "trainer") rows = rows.filter((row) => trainer && row.trainer_id === trainer.id);
  rows = rows.filter((row) => hit([trainerName(data, row.trainer_id), WORK[row.work_format_id], row.issued_on, row.expires_on]));
  return `<h2>Допуски</h2>
    <p class="lead">Допуск тренера к формату: индивид, группа, спорт. Свободное без допуска.</p>
    ${tbl(["Тренер", "Формат", "Выдан", "До"], rows.map((row) => `<tr>
      <td class="name">${esc(trainerName(data, row.trainer_id))}</td>
      <td>${esc(WORK[row.work_format_id] || row.work_format_id)}</td>
      <td>${esc(row.issued_on)}</td><td>${esc(row.expires_on || "бессрочно")}</td>
    </tr>`).join(""))}`;
}

async function journalPage() {
  const params = new URLSearchParams();
  if (filters.action) params.set("action", filters.action);
  if (filters.from) params.set("from", filters.from);
  if (filters.to) params.set("to", filters.to);
  if (q.trim()) params.set("q", q.trim());
  const rows = await load("journal:" + params.toString(), "/journal/?" + params.toString());
  return `<h2>Журнал</h2>
    <p class="lead">Вход и изменения: login, create, update, delete.</p>
    <div class="filters">
      <input id="q" placeholder="поиск по описанию" value="${esc(q)}" />
      <select data-filter="action">
        <option value="">все действия</option>
        ${["login", "create", "update", "delete"].map((action) => `<option value="${action}"${filters.action === action ? " selected" : ""}>${action}</option>`).join("")}
      </select>
      <input data-filter="from" type="date" value="${esc(filters.from)}" />
      <input data-filter="to" type="date" value="${esc(filters.to)}" />
    </div>
    ${tbl(["Время", "Кто", "Действие", "Объект", "Описание", "IP", "Результат"], rows.map((row) => `<tr>
      <td>${esc(row.occurred_at)}</td><td class="name">${esc(row.login)}</td><td>${esc(row.action)}</td>
      <td>${esc(row.object)}</td><td>${esc(row.description)}</td><td>${esc(row.ip_address)}</td>
      <td>${tag(RESULT[row.result] || row.result, row.result === "error")}</td>
    </tr>`).join(""))}`;
}

function statusLink(path, status, bucket, label) {
  return `<a href="#" data-act="status" data-path="${path}" data-status="${status}" data-drop="${bucket}">${label}</a>`;
}
function shell(main) {
  const links = me ? NAV[me.role].map(([id, title]) => `<a href="#" data-p="${id}" class="${page === id ? "on" : ""}">${title}</a>`).join("") : "";
  const who = me ? `<div class="who">${esc(personName(me))}<br>${esc(ROLE[me.role] || me.role)}<br><button type="button" class="ghost" data-out>Выйти</button></div>` : "";
  return `<aside><span class="brand">Бассейн</span><nav>${links}</nav>${who}</aside><main>${main}</main>`;
}
function showMsg(message) {
  let el = document.querySelector("main .msg");
  if (!el) {
    el = document.createElement("p");
    el.className = "msg";
    const lead = document.querySelector("main .lead");
    if (lead) lead.insertAdjacentElement("afterend", el);
    else document.querySelector("main")?.prepend(el);
  }
  el.textContent = explain(message);
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
    return;
  }
  app.innerHTML = shell(`<p class="lead">Загрузка…</p>`);
  let html;
  try {
    html = await pages[page]();
  } catch (error) {
    if (error.status === 401) {
      me = null;
      store = {};
      if (token === renderToken) render();
      return;
    }
    html = `<h2>Ошибка</h2><p class="msg">${esc(explain(error.message))}</p>`;
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
      }});
      drop("users");
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
    showMsg(error.message);
  } finally {
    saving = false;
  }
}

document.querySelectorAll("[data-m]").forEach((button) => {
  button.onclick = () => {
    document.querySelectorAll("[data-m]").forEach((item) => item.classList.remove("on"));
    button.classList.add("on");
    document.body.classList.remove("theme-light", "theme-dark");
    document.body.classList.add(button.dataset.m);
  };
});

document.getElementById("app").addEventListener("click", async (event) => {
  const nav = event.target.closest("a[data-p]");
  if (nav) {
    event.preventDefault();
    page = nav.dataset.p;
    form = null;
    q = "";
    people = "accounts";
    render();
    return;
  }
  const sub = event.target.closest("a[data-s]");
  if (sub) {
    event.preventDefault();
    people = sub.dataset.s;
    form = null;
    q = "";
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
      showMsg(error.message);
    }
    return;
  }
  if (event.target.closest("[data-form=new]")) { form = "new"; render(); return; }
  if (event.target.closest("[data-form=off]")) { form = null; render(); return; }
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
    }
    render();
  } catch (error) {
    showMsg(error.message);
  }
});
document.getElementById("app").addEventListener("input", (event) => {
  if (event.target.id !== "q") return;
  q = event.target.value;
  focusSearch = true;
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => render(), page === "journal" ? 250 : 0);
});
document.getElementById("app").addEventListener("change", (event) => {
  if (event.target.name === "session_id") filterEnrollmentServices();
  if (!event.target.dataset.filter) return;
  filters[event.target.dataset.filter] = event.target.value;
  drop("journal");
  render();
});

api("/me").then((user) => {
  me = user;
  page = NAV[user.role][0][0];
  render();
}).catch(() => render());
