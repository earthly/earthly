package spec

// Earthfile is the AST representation of an Earthfile.
type Earthfile struct {
	Version        *Version        `json:"version,omitempty"`
	BaseRecipe     Block           `json:"baseRecipe"`
	Targets        []Target        `json:"targets,omitempty"`
	UserCommands   []UserCommand   `json:"userCommands,omitempty"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// Target is the AST representation of an Earthfile target.
type Target struct {
	Name           string          `json:"name"`
	Docs           string          `json:"docs,omitempty"`
	Recipe         Block           `json:"recipe"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// UserCommand is the AST representation of an Earthfile user command definition.
type UserCommand struct {
	Name           string          `json:"name"`
	Recipe         Block           `json:"recipe"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// Version is the AST representation of an Earthfile version definition.
type Version struct {
	Args           []string        `json:"args"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// Block is the AST representation of a block of statements.
type Block []Statement

// Statement is the AST representation of an Earthfile statement. Only one field may be
// filled at one time.
type Statement struct {
	Command        *Command        `json:"command,omitempty"`
	With           *WithStatement  `json:"with,omitempty"`
	If             *IfStatement    `json:"if,omitempty"`
	Try            *TryStatement   `json:"try,omitempty"`
	For            *ForStatement   `json:"for,omitempty"`
	Wait           *WaitStatement  `json:"wait,omitempty"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// Command is the AST representation of an Earthfile command.
type Command struct {
	Name           string          `json:"name"`
	Docs           string          `json:"docs,omitempty"`
	Args           []string        `json:"args"`
	ExecMode       bool            `json:"execMode,omitempty"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// WithStatement is the AST representation of a with statement.
type WithStatement struct {
	Command        Command         `json:"command"`
	Body           Block           `json:"body"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// IfStatement is the AST representation of an if statement.
type IfStatement struct {
	Expression     []string        `json:"expression"`
	ExecMode       bool            `json:"execMode,omitempty"`
	IfBody         Block           `json:"ifBody"`
	ElseIf         []ElseIf        `json:"elseIf,omitempty"`
	ElseBody       *Block          `json:"elseBody,omitempty"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// TryStatement is the AST representation of a try statement.
type TryStatement struct {
	TryBody        Block           `json:"tryBody"`
	CatchBody      *Block          `json:"catchBody,omitempty"`
	FinallyBody    *Block          `json:"finallyBody,omitempty"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// ElseIf is the AST representation of an else if clause.
type ElseIf struct {
	Expression     []string        `json:"expression"`
	ExecMode       bool            `json:"execMode,omitempty"`
	Body           Block           `json:"body"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// ForStatement is the AST representation of a for statement.
type ForStatement struct {
	Args           []string        `json:"args"`
	Body           Block           `json:"body"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// WaitStatement is the AST representation of a for statement.
type WaitStatement struct {
	Args           []string        `json:"args"`
	Body           Block           `json:"body"`
	SourceLocation *SourceLocation `json:"sourceLocation,omitempty"`
}

// SourceLocation is an optional reference to the original source code location.
type SourceLocation struct {
	File        string `json:"file,omitempty"`
	StartLine   int    `json:"startLine"`
	StartColumn int    `json:"startColumn"`
	EndLine     int    `json:"endLine"`
	EndColumn   int    `json:"endColumn"`
}
