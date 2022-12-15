package errors

func New(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// The type for these constants is a string. If we need to send an error to the function one level above the call stack,
// We can call errors.New(errors.Example)

const CookieUnfound = "Nu am putut identifica cookie-ul necesar. Asigura-te ca esti inregistrat si " +
	"ai permis utilizarea cookie-urile esentiale."
const CookieValueUnfound = "Hm...Se pare ca nu esti conecat. Te rugam sa te autentifici. "

const SessionExpired = "Sesiunea a expirat. Te rugam sa te autentifici din nou."

const CommentError = "Continut invalid pentru a putea crea un comentariu"
const InsertCommentError = "Nu s-a putut finaliza crearea unui comentariu"

const PostCreateError = "Continut invalid pentru a putea crea o postare."
const UpdatePostError = "Continut invalid pentru a putea actualiza postarea"

const InvalidTitle = "Titlu invalid"

const InvalidUserFields = "Campuri necesare: Nume, Prenume, Email, Parola"

const EmailAlreadyExists = "Email deja inregistrat"
const TokenEmailMismatch = "Token-ul transmis corespunde unei alte adrese de email."
const EmailSendingError = "Nu s-a putut finaliza trimiterea mesajului de confirmare pe email. " +
	"Va rugam sa reincercati procedura de inregistrare a unui cont."
const RequiredEmailPass = "Email-ul si parola sunt necesare pentru a continua"
const InvalidEmailFormat = "Format invalid al email-ului"
const EmailUnfound = "Nu am putut identifica un cont asociat acestei adrese de email."

const InvalidToken = "Token-ul nu este valid. Va rugam reincercati."
const ExpiredToken = "Token-ul transmis prin email a expirat. Va rugam sa reincepeti procedura de inregistrare."

const NotEnoughParameters = "Nu au fost introduci suficienti parametri pentru a finaliza cererea."
const InvalidLoginCredentials = "Email-ul si/sau parola sunt gresite."

const Required = "Acest camp nu poate fi gol."
const PasswordFormat = "Parola trebuie sa contina: " +
	"intre 7 si 49 de caractere, o litera mare, o litera mica si sa aibe doar caractere valide. "
const IdenticalFields = "Campurile trebuie sa fie identice!"
const DifferentFields = "Campurile trebuie sa fie diferite!"
const InvalidPassword = "Parola actuala este incorecta."
