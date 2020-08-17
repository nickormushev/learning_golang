package dictionary

//Dictionary represents a dictionary.
type Dictionary map[string]string

//DictError is a type used for dictionary errors. In essence it is just
//a string but allows for errors to be of a constantType and Immutable
//Structures on the other hand can not be constants and in turn errors.New()
//can't be a const due to it usins a struct
type DictError string

func (d DictError) Error() string {
	return string(d)
}

const (
	//ErrWordNotFound is used when a word has no definition in the dictionary
	ErrWordNotFound = DictError("Word missing from dictionary")
	//ErrAlreadyAdded is thrown when we are attempting to add a definition to a word that already has one
	ErrAlreadyAdded = DictError("Word already added to dictionary")
)

//Search finds a word in a dictionary and returns it's definition
func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]

	if !ok {
		return "", ErrWordNotFound
	}

	return definition, nil
}

//Add a word definition to the dictionary
func (d Dictionary) Add(word, definition string) error {
	_, err := d.Search(word)

	switch err {
	case ErrWordNotFound:
		d[word] = definition
	case nil:
		return ErrAlreadyAdded
	default:
		return err
	}

	return nil
}

//Delete removes a word definition from the dictionary
func (d Dictionary) Delete(word string) {
	delete(d, word)
}

//Update a word definition in the dictionary
func (d Dictionary) Update(word, newDefinition string) error {
	_, err := d.Search(word)

	if err == nil {
		d[word] = newDefinition
	}

	return err
}