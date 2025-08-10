package email

// Purpose: Construct complex objects step-by-step
// Real Use: Document generators (PDF/HTML), complex configuration objects

type Email struct {
	From        string
	To          string
	Subject     string
	Body        string
	Attachments []string
}

type EmailBuilder struct {
	email Email
}

func (eb *EmailBuilder) SetFrom(from string) *EmailBuilder {
	eb.email.From = from
	return eb
}

func (eb *EmailBuilder) SetTo(to string) *EmailBuilder {
	eb.email.To = to
	return eb
}

func (eb *EmailBuilder) SetSubject(subject string) *EmailBuilder {
	eb.email.Subject = subject
	return eb
}

func (eb *EmailBuilder) SetBody(body string) *EmailBuilder {
	eb.email.Body = body
	return eb
}

func (eb *EmailBuilder) AddAttachment(attachment string) *EmailBuilder {
	eb.email.Attachments = append(eb.email.Attachments, attachment)
	return eb
}

func (eb *EmailBuilder) Build() *Email {
	return &eb.email
}

// Usage:
// builder := &EmailBuilder{}
// email := builder.SetFrom("me@example.com").
//          SetTo("you@example.com").
//          AddAttachment("report.pdf").
//          Build()
