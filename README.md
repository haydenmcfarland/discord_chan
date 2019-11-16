# String vs []byte
No reason to force string input when you can use []byte.

# Errors
the error string should never be capitalized nor end with a punctuation per Go standards

# if statement

f, contains := something
if contain {
}

Instead do 

if f, contains := something; contains {

}

# Use iterative algorithms instead of recursive

types - declares structs and possibly some mutators of these structs,
repository - it’s a data storage layer that deals with storing and reading structs,
service - would be the implementation of business logic that wraps repositories,
http, websocket, … - the transport layers, which all invoke the service layer