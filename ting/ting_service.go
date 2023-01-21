package ting

func saveTing(ting Ting) error {
	statement, err := Prepare("insert into ting (audio_url, content, created_at, description, program_id, title, updated_at) values (?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return err
	}

	_, err = statement.Exec(ting.AudioUrl, ting.Content, ting.CreatedAt, ting.Description, ting.ProgramId, ting.Title, ting.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}
