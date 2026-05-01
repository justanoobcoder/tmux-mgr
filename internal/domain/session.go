package domain

type SavedSession struct {
	Name string
}

type DisplaySession struct {
	Name         string
	IsActive     bool
	IsSaved      bool
	IsLastActive bool
	IsLastSaved  bool
}
