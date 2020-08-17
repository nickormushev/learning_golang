package dictionary

import "testing"

func TestDelete(t *testing.T) {
	dictionary := Dictionary{}

	word := "Keyword"
	definition := "Result"

	dictionary.Add(word, definition)

	dictionary.Delete(word)
	_, err := dictionary.Search(word)

	//Word should be deleted
	assertError(t, err, ErrWordNotFound, word)
}

func TestUpdate(t *testing.T) {
	t.Run("Update existing word", func(t *testing.T) {
		dictionary := Dictionary{}
		givenKey := "Keyword"
		value := "Result"
		expectedNewValue := "NewResult"

		dictionary.Add(givenKey, value)

		err := dictionary.Update(givenKey, expectedNewValue)
		assertErrorIsNil(t, err)

		gotValue, err := dictionary.Search("Keyword")

		assertErrorIsNil(t, err)
		assertStrings(t, gotValue, expectedNewValue, givenKey)

	})

	t.Run("Update missing word", func(t *testing.T) {
		dictionary := Dictionary{}
		givenKey := "Keyword"
		givenValue := "Result"

		err := dictionary.Update(givenKey, givenValue)
		assertError(t, err, ErrWordNotFound, givenKey)
	})

}

func TestAdd(t *testing.T) {
	t.Run("Add new word", func(t *testing.T) {

		dictionary := Dictionary{}

		given := "Keyword"
		want := "Result"

		err := dictionary.Add(given, want)
		assertErrorIsNil(t, err)

		got, err := dictionary.Search("Keyword")

		assertErrorIsNil(t, err)
		assertStrings(t, got, want, given)
	})

	t.Run("Attempt to add alredy added word", func(t *testing.T) {
		dictionary := Dictionary{}

		given := "Keyword"
		want := "Result"
		dictionary.Add(given, want)

		err := dictionary.Add(given, want)
		assertError(t, err, ErrAlreadyAdded, given)
	})
}

func TestSearch(t *testing.T) {
	t.Run("Existing word", func(t *testing.T) {
		dictionary := Dictionary{"test": "To validate something"}

		given := "test"
		got, err := dictionary.Search(given)
		want := "To validate something"

		assertErrorIsNil(t, err)
		assertStrings(t, got, want, given)
	})

	t.Run("A word without a definition", func(t *testing.T) {
		dictionary := Dictionary{}

		given := "test"
		_, err := dictionary.Search(given)

		assertError(t, err, ErrWordNotFound, given)
	})

}

func assertError(t *testing.T, got, want error, given string) {
	t.Helper()

	assertErrorNotNil(t, got)

	if got != want {
		t.Errorf("Wanted error: %q, got error: %q, given: %s", want, got, given)
	}
}

func assertStrings(t *testing.T, got, want, given string) {
	t.Helper()

	if got != want {
		t.Errorf("Wanted: %s, got: %s, given: %s", want, got, given)
	}
}

func assertErrorIsNil(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		//Unlike errorf Fatal stops the test from continuing
		t.Fatal("Expected nil but got an error")
	}
}

func assertErrorNotNil(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		//Unlike errorf Fatal stops the test from continuing
		t.Fatal("Expected error but got an nil")
	}
}
