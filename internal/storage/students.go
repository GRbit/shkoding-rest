package storage

type Student struct {
	ID int64
	Name string
	Telegram string
}

func (s *Storage) GetStudent(id int64) (*Student, bool) {
	s.m.RLock()
	defer s.m.RUnlock()

	ret, ok := s.m.Students[id]

	return ret, ok
}

func (s *Storage) GetStudents() ([]*Student) {
	s.m.RLock()
	defer s.m.RUnlock()

	students := make([]*Student, len(s.m.Students), 0)
	for id := range s.m.Students {
		students = append(students, s.m.Students[id])
	}

	return students
}

func (s *Storage) NewStudent(name, tg string) *Student {
	s.m.Lock()
	defer s.m.Unlock()

	s.m.studentsIncrement += 1
	newStudent := &Student{
		ID: s.m.studentsIncrement,
		Name:     name,
		Telegram: tg,
	}

	s.m.Students[s.m.studentsIncrement] = newStudent

	return newStudent
}

func (s *Storage) UpdateStudent(id int64, name, tg string) (*Student, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	student, ok := s.m.Students[id]
	if !ok {
		return nil, false
	}

	if name != "" {
		student.Name = name
	}
	if tg != "" {
		student.Telegram = tg
	}

	//s.m.Students[id] = student

	return student, true
}

func (s *Storage) DeleteStudent(id int64) {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.m.Students, id)

	return
}
