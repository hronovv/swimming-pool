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
const NAMES = { guest: "", client: "Петрова М.Н.", trainer: "Иванов П.А.", admin: "Базан И.А." };
const CLIENT = "Петрова М.Н.";
const TRAINER = "Иванов П.А.";

const TARIFFS = [
  ["1", "Свободное разовое", "свободное", "12.00", "1", "1"],
  ["2", "Свободное абонемент", "свободное", "80.00", "12", "30"],
  ["3", "Индивид разовое", "индивид", "25.00", "1", "1"],
  ["4", "Индивид абонемент", "индивид", "180.00", "8", "30"],
  ["5", "Группа разовое", "группа", "15.00", "1", "1"],
  ["6", "Группа абонемент", "группа", "90.00", "12", "30"],
  ["7", "Спорт-пасс", "спорт", "0.00", "—", "180"],
];
const LANES = [
  ["1", "Дорожка 1", "25", "1.20", "доступна"],
  ["2", "Дорожка 2", "25", "1.40", "доступна"],
  ["3", "Дорожка 3", "25", "1.60", "доступна"],
  ["4", "Дорожка 4", "25", "1.80", "доступна"],
  ["5", "Дорожка 5", "50", "1.80", "доступна"],
  ["6", "Дорожка 6", "50", "2.00", "доступна"],
  ["7", "Дорожка 7", "50", "2.20", "доступна"],
  ["8", "Дорожка 8", "50", "2.20", "закрыта"],
];
const GROUPS = [
  ["Группа 8", "завершена", "18", "21.05.2025 — 17.11.2025"],
  ["Группа 9", "тренировки", "19", "10.06.2025 — 07.12.2025"],
  ["Группа 10", "тренировки", "20", "30.06.2025 — 27.12.2025"],
  ["Группа 11", "тренировки", "10", "20.07.2025 — 16.01.2026"],
  ["Группа 23", "набор", "11", "17.03.2026 — 13.09.2026"],
  ["Группа 24", "набор", "12", "06.04.2026 — 03.10.2026"],
];
const SERVICES = [
  ["1", CLIENT, "Свободное разовое", "—", "активна", "06.01.2025", "07.01.2025"],
  ["586", CLIENT, "Группа разовое", "Группа 6", "активна", "11.01.2025", "12.01.2025"],
  ["1171", CLIENT, "Свободное абонемент", "—", "закрыта", "16.01.2025", "15.02.2025"],
  ["2", "Сидоров П.М.", "Свободное абонемент", "—", "активна", "11.01.2025", "10.02.2025"],
  ["587", "Сидоров П.М.", "Группа абонемент", "Группа 7", "заморожена", "16.01.2025", "15.02.2025"],
  ["3", "Козлова М.", "Индивид разовое", "—", "активна", "16.01.2025", "17.01.2025"],
  ["588", "Козлова М.", "Спорт-пасс", "Группа 28", "активна", "21.01.2025", "20.07.2025"],
];
const PAYMENTS = [
  ["1", CLIENT, "услуга 1", "12.00", "06.01.2025", "карта", "Оплата тарифа Свободное разовое", "выполнен"],
  ["1171", CLIENT, "услуга 1171", "80.00", "16.01.2025", "наличные", "Оплата тарифа Свободное абонемент", "выполнен"],
  ["2", "Сидоров П.М.", "услуга 2", "80.00", "11.01.2025", "перевод", "Оплата тарифа Свободное абонемент", "ожидает"],
  ["3", "Козлова М.", "услуга 3", "25.00", "16.01.2025", "карта", "Оплата тарифа Индивид разовое", "отменён"],
];
const SESSIONS = [
  ["17.11.2025", "08:00–08:45", "свободное", "—", "Дорожка 1", "Свободное плавание", "—", "проведено"],
  ["17.11.2025", "10:00–10:45", "индивид", "—", "Дорожка 3", "Персональная тренировка", TRAINER, "проведено"],
  ["18.11.2025", "09:00–09:45", "группа", "Группа 9", "Дорожка 2", "Групповое занятие", "Волкова А.Н.", "запланировано"],
  ["18.11.2025", "11:00–11:45", "спорт", "Группа 11", "Дорожка 5", "Спортивная тренировка", TRAINER, "запланировано"],
  ["19.11.2025", "08:00–08:45", "индивид", "—", "—", "Персональная тренировка", "Соловьев И.М.", "запланировано"],
];
const ENROLLMENTS = [
  ["17.11.2025 08:00", CLIENT, "Свободное разовое", "да", "Медсправка проверена"],
  ["18.11.2025 09:00", CLIENT, "Группа разовое", "нет", "Опоздал на 10 минут"],
  ["17.11.2025 10:00", "Сидоров П.М.", "Индивид абонемент", "да", "Занятие пройдено полностью"],
  ["18.11.2025 11:00", "Козлова М.", "Спорт-пасс", "нет", "—"],
];
const STANDARDS = [
  [CLIENT, "Волкова А.Н.", "кроль", "25", "02.03.2025", "00:00:26", "сдан", "2"],
  ["Сидоров П.М.", "Соловьев И.М.", "брасс", "50", "03.03.2025", "00:00:27", "сдан", "3"],
  ["Козлова М.", "Васильева А.Е.", "баттерфляй", "100", "04.03.2025", "00:00:28", "сдан", "1"],
  ["Новиков П.П.", "Смирнов И.", "на спине", "200", "05.03.2025", "00:00:29", "сдан", "2"],
  ["Морозова М.И.", "Кузнецова А.И.", "комплекс", "400", "06.03.2025", "—", "не сдан", "3"],
];
const USERS = [
  ["admin", "Базан И.А.", "admin", "—", "—", "активен"],
  ["ivanov.coach10", "Иванов П.А.", "trainer", "+375291234567", "ivanov.coach10@gmail.com", "активен"],
  ["maria.petrov16", CLIENT, "client", "+375293000016", "maria.petrov16@gmail.com", "активен"],
  ["petr.sidorov17", "Сидоров П.М.", "client", "+375333000017", "petr.sidorov17@mail.ru", "активен"],
  ["maria.kozlov18", "Козлова М.", "client", "+375443000018", "maria.kozlov18@yandex.by", "заблокирован"],
];
const CLIENTS = [
  [CLIENT, "MP1159568", "25.10.1969", "г. Минск, ул. Якуба Коласа, д. 61, кв. 23", "нет"],
  ["Сидоров П.М.", "BM1169541", "01.12.1969", "г. Минск, ул. Кирова, д. 74, кв. 30", "есть"],
  ["Козлова М.", "AB1179514", "07.01.1970", "г. Минск, ул. Максима Богдановича, д. 87, кв. 37", "есть"],
  ["Новиков П.П.", "HB1189487", "13.02.1970", "г. Минск, ул. Сурганова, д. 100, кв. 44", "есть"],
];
const TRAINERS = [
  ["Волкова А.Н.", "инструктор", "нет", "6", "кроль", "индивид, группа, спорт"],
  ["Соловьев И.М.", "тренер", "3 разряд", "7", "брасс", "индивид, группа, спорт"],
  ["Васильева А.Е.", "старший тренер", "2 разряд", "8", "баттерфляй", "индивид, группа, спорт"],
  [TRAINER, "инструктор", "1 разряд", "15", "комплекс", "индивид"],
];
const PERMITS = [
  [TRAINER, "индивид", "01.01.2024", "бессрочно"],
  ["Волкова А.Н.", "индивид", "01.01.2024", "бессрочно"],
  ["Волкова А.Н.", "группа", "01.01.2024", "бессрочно"],
  ["Волкова А.Н.", "спорт", "01.01.2024", "бессрочно"],
  ["Смирнов И.", "группа", "01.01.2024", "30.06.2027"],
];
const JOURNAL = [
  ["11.11 10:32", "admin", "login", "users", "Вход в систему", "успех", "192.0.2.1"],
  ["11.11 10:33", "admin", "create", "payments", "Платёж 50.00", "успех", "192.0.2.1"],
  ["11.11 10:34", "admin", "update", "services", "Услуга 587 заморожена", "успех", "192.0.2.1"],
  ["10.11 18:02", "maria.petrov16", "login", "users", "Вход в систему", "успех", "192.0.2.16"],
  ["10.11 18:05", "maria.petrov16", "login", "users", "Неверный пароль", "ошибка", "192.0.2.16"],
];

let role = "admin", page = "services", q = "", form = null, people = "accounts";

const tag = (t, warn) => `<span class="tag${warn ? " warn" : ""}">${t}</span>`;
const esc = s => String(s).replace(/"/g, "&quot;");
function match(row) {
  return !q.trim() || row.join(" ").toLowerCase().includes(q.toLowerCase());
}
function tbl(heads, rowsHtml) {
  return `<div class="table-wrap"><table><thead><tr>${heads.map(h => `<th>${h}</th>`).join("")}</tr></thead><tbody>${rowsHtml}</tbody></table></div>`;
}
function formBtns() {
  return `<div class="wide"><button type="button" class="btn">Сохранить</button> <button type="button" class="btn-line" data-form="off">Отмена</button></div>`;
}
function navHtml() {
  return NAV[role].map(([id, t]) => `<a href="#" data-p="${id}" class="${page === id ? "on" : ""}">${t}</a>`).join("");
}
function searchBar(ph) {
  const add = (ADDABLE[role] || []).includes(page);
  return `<div class="toolbar">
    <input id="q" placeholder="${ph}" value="${esc(q)}" />
    ${add ? `<button type="button" class="btn" data-form="new">Добавить</button>` : ""}
  </div>`;
}

function loginPage() {
  return `<div class="login-page">
    <h1>Вход</h1>
    <p class="lead">Логин и пароль учётной записи бассейна.</p>
    <div class="field"><label>Логин</label><input value="maria.petrov16" /></div>
    <div class="field"><label>Пароль</label><input type="password" value="password" /></div>
    <button type="button" class="btn" data-enter>Войти</button>
  </div>`;
}
function servicesPage() {
  const rows = (role === "client" ? SERVICES.filter(r => r[1] === CLIENT) : SERVICES).filter(match);
  return `<h2>${role === "client" ? "Мои услуги" : "Услуги"}</h2>
    <p class="lead">${role === "client" ? "Оформленные тарифы." : "Оформление услуги клиенту, смена статуса."}</p>
    ${searchBar("Поиск: клиент, тариф, группа")}
    ${form === "new" && role === "admin" ? `<div class="form-grid">
      <div><label>Клиент</label><select><option>${CLIENT}</option><option>Сидоров П.М.</option></select></div>
      <div><label>Тариф</label><select>${TARIFFS.map(t => `<option>${t[1]}</option>`).join("")}</select></div>
      <div><label>Группа</label><select><option value="">—</option>${GROUPS.map(g => `<option>${g[0]}</option>`).join("")}</select></div>
      <div><label>Начало</label><input type="date" value="2025-11-17" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["id", "Клиент", "Тариф", "Группа", "Статус", "Срок"].concat(role === "admin" ? [""] : []),
      rows.map(r => `<tr>
        <td>${r[0]}</td><td class="name">${r[1]}</td><td>${r[2]}</td><td>${r[3]}</td>
        <td>${tag(r[4], r[4] === "заморожена")}</td><td>${r[5]} — ${r[6]}</td>
        ${role === "admin" ? `<td class="row-actions"><a href="#">заморозить</a><a href="#">закрыть</a></td>` : ""}
      </tr>`).join(""))}`;
}
function paymentsPage() {
  const rows = (role === "client" ? PAYMENTS.filter(r => r[1] === CLIENT) : PAYMENTS).filter(match);
  return `<h2>Платежи</h2>
    <p class="lead">${role === "client" ? "Оплаты по своим услугам." : "Внесение оплаты по услуге. Спорт-пасс без платежа."}</p>
    ${searchBar("Поиск: клиент, сумма, назначение")}
    ${form === "new" && role === "admin" ? `<div class="form-grid">
      <div><label>Услуга</label><select><option>1 — ${CLIENT}</option></select></div>
      <div><label>Сумма</label><input value="12.00" /></div>
      <div><label>Способ</label><select><option>карта</option><option>наличные</option><option>перевод</option></select></div>
      <div><label>Назначение</label><input value="Оплата тарифа" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Клиент", "Услуга", "Сумма", "Дата", "Способ", "Назначение", "Статус"],
      rows.map(r => `<tr>
        <td class="name">${r[1]}</td><td>${r[2]}</td><td class="money">${r[3]}</td>
        <td>${r[4]}</td><td>${r[5]}</td><td>${r[6]}</td>
        <td>${tag(r[7], r[7] !== "выполнен")}</td>
      </tr>`).join(""))}`;
}
function sessionsPage() {
  let rows = SESSIONS;
  if (role === "trainer") rows = SESSIONS.filter(r => r[6] === TRAINER);
  if (role === "client") rows = SESSIONS.filter(r => r[2] === "свободное" || r[3] === "Группа 9");
  rows = rows.filter(match);
  return `<h2>${role === "trainer" ? "Расписание" : "Занятия"}</h2>
    <p class="lead">${role === "trainer" ? "Свои занятия." : role === "client" ? "Ближайшие занятия." : "Сессии: формат, группа, дорожка, тренер."}</p>
    ${searchBar("Поиск: дата, формат, дорожка")}
    ${form === "new" && role === "admin" ? `<div class="form-grid">
      <div><label>Дата</label><input type="date" value="2025-11-19" /></div>
      <div><label>Формат</label><select><option>свободное</option><option>индивид</option><option>группа</option><option>спорт</option></select></div>
      <div><label>Дорожка</label><select>${LANES.map(l => `<option>${l[1]}</option>`).join("")}</select></div>
      <div><label>Группа</label><select><option value="">—</option>${GROUPS.map(g => `<option>${g[0]}</option>`).join("")}</select></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Дата", "Время", "Формат", "Группа", "Дорожка", "Тема", "Тренер", "Статус"],
      rows.map(r => `<tr>
        <td class="name">${r[0]}</td><td>${r[1]}</td><td>${r[2]}</td><td>${r[3]}</td>
        <td>${r[4]}</td><td>${r[5]}</td><td>${r[6]}</td><td>${tag(r[7])}</td>
      </tr>`).join(""))}`;
}
function enrollmentsPage() {
  let rows = ENROLLMENTS;
  if (role === "client") rows = ENROLLMENTS.filter(r => r[1] === CLIENT);
  rows = rows.filter(match);
  return `<h2>Явка</h2>
    <p class="lead">Запись услуги на занятие и отметка посещения.</p>
    ${searchBar("Поиск: клиент, занятие")}
    ${form === "new" && role === "admin" ? `<div class="form-grid">
      <div><label>Занятие</label><select><option>18.11.2025 09:00 — Группа 9</option></select></div>
      <div><label>Услуга</label><select><option>586 — ${CLIENT}</option></select></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Занятие", "Клиент", "Тариф", "Явка", "Комментарий"].concat(role !== "client" ? [""] : []),
      rows.map(r => `<tr>
        <td class="name">${r[0]}</td><td>${r[1]}</td><td>${r[2]}</td>
        <td>${tag(r[3], r[3] === "нет")}</td><td>${r[4]}</td>
        ${role !== "client" ? `<td class="row-actions"><a href="#">${r[3] === "да" ? "снять явку" : "отметить"}</a></td>` : ""}
      </tr>`).join(""))}`;
}
function standardsPage() {
  let rows = STANDARDS;
  if (role === "client") rows = STANDARDS.filter(r => r[0] === CLIENT);
  rows = rows.filter(match);
  return `<h2>Нормативы</h2>
    <p class="lead">${role === "client" ? "Свои попытки." : "Стиль, дистанция, время, сдан / не сдан."}</p>
    ${searchBar("Поиск: клиент, стиль")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Клиент</label><select><option>${CLIENT}</option></select></div>
      <div><label>Стиль</label><select><option>кроль</option><option>брасс</option><option>баттерфляй</option><option>на спине</option><option>комплекс</option></select></div>
      <div><label>Дистанция, м</label><select><option>25</option><option>50</option><option>100</option><option>200</option><option>400</option></select></div>
      <div><label>Результат</label><input placeholder="00:00:26" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Клиент", "Тренер", "Стиль", "м", "Дата", "Время", "Итог", "Попытка"],
      rows.map(r => `<tr>
        <td class="name">${r[0]}</td><td>${r[1]}</td><td>${r[2]}</td><td>${r[3]}</td>
        <td>${r[4]}</td><td class="money">${r[5]}</td>
        <td>${tag(r[6], r[6] === "не сдан")}</td><td>${r[7]}</td>
      </tr>`).join(""))}`;
}
function usersPage() {
  const tabs = [["accounts", "Счета"], ["clients", "Клиенты"], ["trainers", "Тренеры"]];
  const sub = `<div class="subnav">${tabs.map(([id, t]) => `<a href="#" data-s="${id}" class="${people === id ? "on" : ""}">${t}</a>`).join("")}</div>`;
  if (people === "clients") {
    const rows = CLIENTS.filter(match);
    return `<h2>Пользователи</h2>
      <p class="lead">Карточка клиента: паспорт, дата рождения, медсправка.</p>
      ${sub}${searchBar("Поиск: ФИО, паспорт")}
      ${tbl(["ФИО", "Паспорт", "Рождение", "Адрес", "Медсправка"],
        rows.map(r => `<tr>
          <td class="name">${r[0]}</td><td>${r[1]}</td><td>${r[2]}</td><td>${r[3]}</td>
          <td>${tag(r[4], r[4] === "нет")}</td>
        </tr>`).join(""))}`;
  }
  if (people === "trainers") {
    const rows = TRAINERS.filter(match);
    return `<h2>Пользователи</h2>
      <p class="lead">Карточка тренера и допуски к форматам.</p>
      ${sub}${searchBar("Поиск: ФИО, квалификация")}
      ${tbl(["ФИО", "Квалификация", "Разряд", "Стаж", "Специализация", "Допуски"],
        rows.map(r => `<tr>
          <td class="name">${r[0]}</td><td>${r[1]}</td><td>${r[2]}</td>
          <td>${r[3]}</td><td>${r[4]}</td><td>${r[5]}</td>
        </tr>`).join(""))}`;
  }
  const rows = USERS.filter(match);
  return `<h2>Пользователи</h2>
    <p class="lead">Учётные записи: login, роль, блокировка.</p>
    ${sub}${searchBar("Поиск: логин, фамилия")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Логин</label><input /></div>
      <div><label>Роль</label><select><option>client</option><option>trainer</option><option>admin</option></select></div>
      <div><label>Фамилия</label><input /></div>
      <div><label>Имя</label><input /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Логин", "ФИО", "Роль", "Телефон", "Почта", "Статус", ""],
      rows.map(r => `<tr>
        <td>${r[0]}</td><td class="name">${r[1]}</td><td>${r[2]}</td>
        <td>${r[3]}</td><td>${r[4]}</td>
        <td>${tag(r[5], r[5] === "заблокирован")}</td>
        <td class="row-actions"><a href="#">${r[5] === "активен" ? "заблокировать" : "активировать"}</a></td>
      </tr>`).join(""))}`;
}
function tariffsPage() {
  const rows = TARIFFS.filter(match);
  return `<h2>Тарифы</h2>
    <p class="lead">Справочник: формат, цена, лимит посещений, срок.</p>
    ${searchBar("Поиск: название, формат")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Название</label><input /></div>
      <div><label>Формат</label><select><option>свободное</option><option>индивид</option><option>группа</option><option>спорт</option></select></div>
      <div><label>Цена</label><input value="0.00" /></div>
      <div><label>Дней</label><input value="30" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["id", "Название", "Формат", "Цена", "Визитов", "Дней"],
      rows.map(r => `<tr>
        <td>${r[0]}</td><td class="name">${r[1]}</td><td>${r[2]}</td>
        <td class="money">${r[3]}</td><td>${r[4]}</td><td>${r[5]}</td>
      </tr>`).join(""))}`;
}
function lanesPage() {
  const rows = LANES.filter(match);
  return `<h2>Дорожки</h2>
    <p class="lead">Длина 25 или 50 м, глубина, доступна / закрыта.</p>
    ${searchBar("Поиск: название, статус")}
    ${form === "new" ? `<div class="form-grid">
      <div><label>Название</label><input value="Дорожка 9" /></div>
      <div><label>Длина, м</label><select><option>25</option><option>50</option></select></div>
      <div><label>Глубина, м</label><input value="1.80" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["id", "Название", "Длина", "Глубина", "Статус"],
      rows.map(r => `<tr>
        <td>${r[0]}</td><td class="name">${r[1]}</td><td>${r[2]} м</td>
        <td>${r[3]} м</td><td>${tag(r[4], r[4] === "закрыта")}</td>
      </tr>`).join(""))}`;
}
function groupsPage() {
  const rows = GROUPS.filter(match);
  return `<h2>Группы</h2>
    <p class="lead">Набор, тренировки, завершена. Вместимость и срок.</p>
    ${searchBar("Поиск: название, статус")}
    ${form === "new" && role === "admin" ? `<div class="form-grid">
      <div><label>Название</label><input /></div>
      <div><label>Вместимость</label><input value="12" /></div>
      <div><label>Начало</label><input type="date" value="2026-03-01" /></div>
      <div><label>Конец</label><input type="date" value="2026-08-28" /></div>
      ${formBtns()}
    </div>` : ""}
    ${tbl(["Название", "Статус", "Мест", "Срок"],
      rows.map(r => `<tr>
        <td class="name">${r[0]}</td><td>${tag(r[1], r[1] === "набор")}</td>
        <td class="money">${r[2]}</td><td>${r[3]}</td>
      </tr>`).join(""))}`;
}
function permitsPage() {
  const rows = (role === "trainer" ? PERMITS.filter(r => r[0] === TRAINER) : PERMITS).filter(match);
  return `<h2>Допуски</h2>
    <p class="lead">Допуск тренера к формату: индивид, группа, спорт. Свободное без допуска.</p>
    ${tbl(["Тренер", "Формат", "Выдан", "До"],
      rows.map(r => `<tr>
        <td class="name">${r[0]}</td><td>${r[1]}</td><td>${r[2]}</td><td>${r[3]}</td>
      </tr>`).join(""))}`;
}
function journalPage() {
  const rows = JOURNAL.filter(match);
  return `<h2>Журнал</h2>
    <p class="lead">Вход и изменения: login, create, update, delete.</p>
    <div class="filters">
      <input id="q" placeholder="поиск по описанию" value="${esc(q)}" />
      <select><option>все действия</option><option>login</option><option>create</option><option>update</option><option>delete</option></select>
      <input type="date" value="2025-11-01" />
      <input type="date" value="2025-11-17" />
    </div>
    ${tbl(["Время", "Кто", "Действие", "Объект", "Описание", "IP", "Результат"],
      rows.map(r => `<tr>
        <td>${r[0]}</td><td class="name">${r[1]}</td><td>${r[2]}</td><td>${r[3]}</td>
        <td>${r[4]}</td><td>${r[6]}</td>
        <td>${tag(r[5], r[5] === "ошибка")}</td>
      </tr>`).join(""))}`;
}
function body() {
  if (role === "guest") return loginPage();
  return {
    services: servicesPage, payments: paymentsPage, sessions: sessionsPage,
    enrollments: enrollmentsPage, standards: standardsPage, users: usersPage,
    tariffs: tariffsPage, lanes: lanesPage, groups: groupsPage,
    permits: permitsPage, journal: journalPage,
  }[page]();
}
function render() {
  const app = document.getElementById("app");
  const name = NAMES[role];
  const login = role === "guest";
  app.className = "layout" + (login ? " login-full" : "");
  app.innerHTML = `<aside>
      <span class="brand">Бассейн</span>
      <nav>${navHtml()}</nav>
      ${name ? `<div class="who">${name}<br><button type="button" class="ghost" data-out>Выйти</button></div>` : ""}
    </aside><main>${body()}</main>`;
}

function setRoleButtons(current) {
  document.querySelectorAll("[data-r]").forEach(x => x.classList.toggle("on", x.dataset.r === current));
}

document.querySelectorAll("[data-m]").forEach(b => {
  b.onclick = () => {
    document.querySelectorAll("[data-m]").forEach(x => x.classList.remove("on"));
    b.classList.add("on");
    document.body.classList.remove("theme-light", "theme-dark");
    document.body.classList.add(b.dataset.m);
  };
});
document.querySelectorAll("[data-r]").forEach(b => {
  b.onclick = () => {
    setRoleButtons(b.dataset.r);
    role = b.dataset.r;
    form = null;
    q = "";
    people = "accounts";
    page = role === "guest" ? "login" : NAV[role][0][0];
    render();
  };
});
document.getElementById("app").addEventListener("click", e => {
  const a = e.target.closest("a[data-p]");
  if (a) { e.preventDefault(); page = a.dataset.p; form = null; q = ""; people = "accounts"; render(); }
  const s = e.target.closest("a[data-s]");
  if (s) { e.preventDefault(); people = s.dataset.s; form = null; q = ""; render(); }
  if (e.target.closest("[data-out]")) {
    role = "guest";
    page = "login";
    setRoleButtons("guest");
    render();
  }
  if (e.target.closest("[data-enter]")) {
    role = "client";
    page = "services";
    setRoleButtons("client");
    render();
  }
  if (e.target.closest("[data-form=new]")) { form = "new"; render(); }
  if (e.target.closest("[data-form=off]")) { form = null; render(); }
});
document.getElementById("app").addEventListener("input", e => {
  if (e.target.id === "q") {
    q = e.target.value;
    render();
    const el = document.getElementById("q");
    if (el) { el.focus(); el.setSelectionRange(q.length, q.length); }
  }
});
render();
