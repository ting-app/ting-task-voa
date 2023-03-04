package ting

import (
	"context"
	"log"
)

func saveTing(ting Ting, tagId int) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, "insert into ting (audio_url, content, created_at, description, program_id, title, updated_at) values (?, ?, ?, ?, ?, ?, ?)", ting.AudioUrl, ting.Content, ting.CreatedAt, ting.Description, ting.ProgramId, ting.Title, ting.UpdatedAt)

	if err != nil {
		tx.Rollback()

		return err
	}

	tingId, err := result.LastInsertId()

	if err != nil {
		tx.Rollback()

		return err
	}

	log.Printf("ting saved, id is %v\n", tingId)

	result, err = tx.ExecContext(ctx, "insert into ting_tag (ting_id, tag_id) values (?, ?)", tingId, tagId)

	if err != nil {
		tx.Rollback()

		return err
	}

	err = tx.Commit()

	return err
}
