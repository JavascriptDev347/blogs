# RBAC vs ABAC — Go'da amaliy taqqoslash

Ruxsatlarni boshqarishning ikkita eng keng tarqalgan modeli — **RBAC** (Role-Based Access Control) va **ABAC** (Attribute-Based Access Control) — o'rtasidagi farqni bitta ishlaydigan misolda ko'rsatuvchi kichik loyiha.

Asosiy savol: *"Seller mahsulotni tahrirlay oladimi?"* — RBAC bu savolga **ha** deb javob beradi, ABAC esa **"qaysi mahsulotni?"** deb so'raydi. Butun farq ana shu qo'shimcha savolda.

---

## Mundarija

- [Modellar haqida qisqacha](#modellar-haqida-qisqacha)
- [Loyiha tuzilishi](#loyiha-tuzilishi)
- [Domen modeli](#domen-modeli)
- [RBAC implementatsiyasi](#rbac-implementatsiyasi)
- [ABAC implementatsiyasi](#abac-implementatsiyasi)
- [Stsenariylar va natijalar](#stsenariylar-va-natijalar)
- [Ishga tushirish](#ishga-tushirish)
- [Xulosa](#xulosa)
- [Keyingi qadamlar](#keyingi-qadamlar)

---

## Modellar haqida qisqacha

| | **RBAC** | **ABAC** |
|---|---|---|
| **Qaror asosi** | Foydalanuvchining **roli** | Subyekt, obyekt, amal va muhit **atributlari** |
| **Savol** | "Bu rolda bunday amal bormi?" | "Aynan shu foydalanuvchi aynan shu obyektga shu sharoitda nima qila oladi?" |
| **Ma'lumot tuzilmasi** | `map[Role][]Action` | Policy'lar (funksiyalar) to'plami |
| **Kuchli tomoni** | Sodda, tushunarli, tez | Moslashuvchan, kontekstni hisobga oladi |
| **Zaif tomoni** | Obyektni ko'rmaydi — "o'z mahsuloti"ni ajrata olmaydi | Murakkabroq, debug qilish qiyinroq |
| **Qachon yetarli** | Ichki admin panellar, oddiy CRUD | Marketplace, multi-tenant, moliyaviy tizimlar |

> **Muhim:** ABAC RBAC'ning o'rnini bosmaydi — odatda uning **ustiga** qo'yiladi. Avval rol bo'yicha filtrlanadi, keyin atributlar bo'yicha aniqlashtiriladi.

---

## Loyiha tuzilishi

```
rbac-vs-abac/
├── go.mod        # module rbac-vs-abac
├── domain.go     # Umumiy tiplar: Role, Action, User, Product, Environment
├── rbac.go       # RBAC: rol → ruxsatlar jadvali
├── abac.go       # ABAC: Policy tipi, OwnershipPolicy, Engine
├── main.go       # 4 ta stsenariy va ularni taqqoslash
└── README.md
```

---

## Domen modeli

`domain.go` — ikkala model ham foydalanadigan umumiy tiplar.

### Rollar

```go
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
	RoleSeller   Role = "seller"
)
```

### Amallar

```go
type Action string

const (
	ActionCreateProduct Action = "create_product"
	ActionUpdateProduct Action = "update_product"
	ActionDeleteProduct Action = "delete_product"
	ActionViewProduct   Action = "view_product"
)
```

### Subyekt, obyekt va muhit

```go
type User struct {          // subyekt — kim?
	ID   string
	Role Role
	Name string
}

type Product struct {       // obyekt — nimaga?
	ID       string
	OwnerID  string          // <- ABAC uchun kalit atribut
	Name     string
	Category string
}

type Environment struct {   // muhit — qanday sharoitda?
	CurrentTime time.Time
	IPAddress   string
}
```

`Product.OwnerID` — butun misolning yuragi. RBAC bu maydonni umuman ko'rmaydi, ABAC esa aynan shunga qarab qaror qabul qiladi.

---

## RBAC implementatsiyasi

`rbac.go` — bor-yo'g'i bitta jadval va bitta funksiya.

```go
var rolePermissions = map[Role][]Action{
	RoleAdmin:    {ActionCreateProduct, ActionDeleteProduct, ActionUpdateProduct, ActionViewProduct},
	RoleCustomer: {ActionViewProduct},
	RoleSeller:   {ActionCreateProduct, ActionUpdateProduct, ActionViewProduct},
}

func CanRBAC(user User, action Action) bool {
	permissions, ok := rolePermissions[user.Role]
	if !ok {
		return false
	}
	for _, perm := range permissions {
		if perm == action {
			return true
		}
	}
	return false
}
```

### Ruxsatlar matritsasi

| Amal | `admin` | `seller` | `customer` |
|---|:---:|:---:|:---:|
| `create_product` | ✅ | ✅ | ❌ |
| `view_product` | ✅ | ✅ | ✅ |
| `update_product` | ✅ | ✅ | ❌ |
| `delete_product` | ✅ | ❌ | ❌ |

**E'tibor bering:** `CanRBAC` imzosida `Product` umuman yo'q:

```go
func CanRBAC(user User, action Action) bool
```

Ya'ni RBAC tabiatan *qaysi* mahsulot ekanini bila olmaydi. Bu uning kamchiligi emas, dizayni shunday.

---

## ABAC implementatsiyasi

`abac.go` — policy'lar to'plami va ularni ketma-ket tekshiruvchi engine.

### Policy tipi

Har bir policy — to'rtala atributni (subyekt, obyekt, amal, muhit) qabul qiluvchi oddiy funksiya:

```go
type Policy func(user User, product Product, action Action, env Environment) bool
```

### Ownership policy

```go
func OwnershipPolicy(user User, product Product, action Action, env Environment) bool {
	if action != ActionUpdateProduct && action != ActionDeleteProduct {
		return true // bu policy faqat update/delete ga tegishli, boshqasiga aralashmaydi
	}
	if user.Role == RoleAdmin {
		return true // admin har doim o'tadi
	}
	return product.OwnerID == user.ID
}
```

Uch qatorlik mantiq, lekin uchta muhim g'oya bor:

1. **Tegishli bo'lmasa — aralashmaydi.** Policy o'z doirasidan tashqaridagi amallarga `true` qaytaradi, ya'ni "menda e'tiroz yo'q" deydi.
2. **Admin uchun istisno.** Rol atributi ham policy ichida ishlatilishi mumkin — ABAC rolni inkor etmaydi, uni shunchaki *bitta atribut* sifatida ko'radi.
3. **Asosiy qaror** — `product.OwnerID == user.ID`. Aynan shu qator RBAC hech qachon ayta olmaydigan gapni aytadi.

### Engine

```go
type Engine struct {
	Policies []Policy
}

func (e Engine) CanABAC(user User, product Product, action Action, env Environment) bool {
	for _, policy := range e.Policies {
		if !policy(user, product, action, env) {
			return false
		}
	}
	return true
}
```

Engine **AND** mantiqida ishlaydi: **barcha** policy'lar rozi bo'lsagina ruxsat beriladi. Bittasi `false` qaytarsa — darhol rad etiladi (fail-closed). Yangi qoida qo'shish uchun mavjud kodni o'zgartirish shart emas, faqat ro'yxatga yangi funksiya qo'shiladi:

```go
engine := Engine{
	Policies: []Policy{OwnershipPolicy},
}
```

---

## Stsenariylar va natijalar

`main.go` bitta mahsulot (egasi — `seller`, ID = `"3"`) va to'rtta foydalanuvchi ustida bir xil savolni ikkala modelga beradi.

| # | Stsenariy | RBAC | ABAC | Izoh |
|:-:|---|:---:|:---:|---|
| 1 | `seller2` **boshqa** sellerning mahsulotini tahrirlaydi | ✅ `true` | ❌ `false` | **Eng muhim qator.** RBAC "seller update qila oladi" deydi va to'xtaydi. ABAC `OwnerID != user.ID` ekanini ko'rib rad etadi. |
| 2 | `seller` **o'z** mahsulotini tahrirlaydi | ✅ `true` | ✅ `true` | Ikkala model ham rozi — bu to'g'ri xatti-harakat. |
| 3 | `admin` istalgan mahsulotni tahrirlaydi | ✅ `true` | ✅ `true` | ABAC'da admin uchun aniq istisno yozilgan. |
| 4 | `customer` mahsulotni tahrirlaydi | ❌ `false` | ❌ `false` | RBAC rol darajasida to'sadi, ABAC ownership darajasida. |

### Haqiqiy chiqish

```
=== Stsenariy 1: seller2 BOSHQA sellerning mahsulotini tahrirlamoqchi ===
RBAC natija: true
ABAC natija: false

=== Stsenariy 2: seller1 O'Z mahsulotini tahrirlamoqchi ===
RBAC natija: true
ABAC natija: true

=== Stsenariy 3: admin istalgan mahsulotni tahrirlamoqchi ===
RBAC natija: true
ABAC natija: true

=== Stsenariy 4: customer mahsulotni tahrirlamoqchi ===
RBAC natija: false
ABAC natija: false
```

**1-stsenariy** — butun loyihaning sababi. Faqat RBAC'ga tayangan marketplace'da har qanday sotuvchi boshqa sotuvchining mahsulotini tahrirlay oladi. Bu jiddiy xavfsizlik teshigi va u sizning kodingizda emas, **modelni tanlashda** yuzaga keladi.

---

## Ishga tushirish

```bash
cd rbac-vs-abac
go run .
```

Talab: Go 1.21+ (loyihada tashqi bog'liqliklar yo'q).

---

## Xulosa

- **RBAC** savolga *"kim?"* darajasida javob beradi — rol va amal. Sodda, o'qilishi oson, aksariyat ichki tizimlar uchun yetarli.
- **ABAC** savolga *"kim, nimaga, qachon, qanday sharoitda?"* darajasida javob beradi. Obyekt atributlarini ko'rgani uchun "o'z resursi" tushunchasini ifodalay oladi.
- Amalda ular **birga** ishlatiladi: RBAC keng filtr sifatida, ABAC esa aniqlashtiruvchi qatlam sifatida.
- ABAC'ning kuchi policy'larning **kompozitsiyasida**: har bir qoida mustaqil, testlanadigan funksiya bo'lib qoladi, ular esa AND mantiqida birlashadi.

Amaliy qoida: *agar ruxsat qarori resursning o'ziga bog'liq bo'lsa — faqat RBAC yetarli emas.*

---

## Keyingi qadamlar

Loyihani kengaytirish uchun tabiiy yo'nalishlar:

- **`RolePolicy`** — RBAC'ni engine ichiga policy sifatida qo'shish, shunda ikkala qatlam bitta joyda birlashadi.
- **`TimePolicy`** — `Environment.CurrentTime`dan foydalanib, tahrirlashni faqat ish vaqtida ruxsat etish.
- **`IPPolicy`** — `Environment.IPAddress` bo'yicha admin amallarini ichki tarmoq bilan cheklash.
- **Qaror sababi** — `bool` o'rniga `(bool, string)` qaytarish, ya'ni "nega rad etildi" ni ham bilish (audit log uchun juda foydali).
- **Testlar** — har bir policy uchun table-driven testlar; policy'lar toza funksiya bo'lgani uchun testlash juda oson.
