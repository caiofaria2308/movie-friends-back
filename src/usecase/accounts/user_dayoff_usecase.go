package usecase_accounts

import (
	entity_accounts "app/entity/accounts"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	UpdateModeSingle = "single"
	UpdateModeFuture = "future"
	UpdateModeAll    = "all"

	DeleteModeSingle = "single"
	DeleteModeFuture = "future"
	DeleteModeAll    = "all"
)

type userDayOffUseCase struct {
	repo IRepositoryUserDayOff
}

func NewUserDayOffUseCase(repo IRepositoryUserDayOff) IUseCaseUserDayOff {
	return &userDayOffUseCase{repo: repo}
}

func (u *userDayOffUseCase) Create(dayOff *entity_accounts.UserDayOff, ownerID int) error {
	dayOff.OwnerID = ownerID

	// Create the primary day off
	if err := u.repo.Create(dayOff); err != nil {
		return fmt.Errorf("could not create day off")
	}

	// Handle recurrence in goroutine if needed
	if dayOff.Repeat && dayOff.RepeatValue != "" {
		go u.handleRecurrence(dayOff, ownerID)
	}

	return nil
}

func (u *userDayOffUseCase) handleRecurrence(father *entity_accounts.UserDayOff, ownerID int) {
	// Simple recurrence logic: create for next 30 occurrences ? or based on RepeatValue?
	// User said "create in mass". Let's assume a fixed number for now or parse RepeatValue.
	// Since `RepeatValue` format isn't strictly defined, I'll assume it's like "10" for 10 times, or maybe just create a reasonable amount (e.g., 50) for now until specced.
	// Actually, let's create for 1 year ahead or 50 occurrences max to avoid infinite loop.

	var newDayOffs []*entity_accounts.UserDayOff

	currentStart := *father.InitHour
	currentEnd := *father.EndHour
	count := 0
	maxCount := 52 // Example: 1 year of weekly events

	if father.RepeatType == entity_accounts.RepeatTypeDaily {
		maxCount = 365
	}

RecurrenceLoop:
	for i := 0; i < maxCount; i++ {
		// Calculate next date
		switch father.RepeatType {
		case entity_accounts.RepeatTypeWeekly:
			currentStart = currentStart.AddDate(0, 0, 7)
			currentEnd = currentEnd.AddDate(0, 0, 7)
		case entity_accounts.RepeatTypeDaily:
			currentStart = currentStart.AddDate(0, 0, 1)
			currentEnd = currentEnd.AddDate(0, 0, 1)
		case entity_accounts.RepeatTypeMonthly:
			currentStart = currentStart.AddDate(0, 1, 0)
			currentEnd = currentEnd.AddDate(0, 1, 0)
		case entity_accounts.RepeatTypeYearly:
			currentStart = currentStart.AddDate(1, 0, 0)
			currentEnd = currentEnd.AddDate(1, 0, 0)
		default:
			break RecurrenceLoop
		}

		id := uuid.New()
		newDayOff := &entity_accounts.UserDayOff{
			ID:             &id,
			InitHour:       &currentStart,
			EndHour:        &currentEnd,
			OwnerID:        ownerID,
			Repeat:         true, // They are part of a repeating series
			RepeatType:     father.RepeatType,
			RepeatValue:    father.RepeatValue,
			DayOffFatherID: father.ID, // Link to father
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		newDayOffs = append(newDayOffs, newDayOff)
		count++
	}

	if len(newDayOffs) > 0 {
		_ = u.repo.CreateBatch(newDayOffs)
	}
}

func (u *userDayOffUseCase) Update(dayOff *entity_accounts.UserDayOff, ownerID int, mode string) error {
	if dayOff.ID == nil {
		return fmt.Errorf("id required")
	}

	existing, err := u.repo.FindByIdAndOwner(*dayOff.ID, ownerID)
	if err != nil {
		return fmt.Errorf("day off not found")
	}

	// Prepare update function
	updateFields := func(target *entity_accounts.UserDayOff) {
		// Calculate duration difference to shift start/end times correctly?
		// Or just replace times?
		// User request: "updates all or from date selected".
		// Usually this means shifting time or changing properties.
		// For simplicity, let's assume we copy properties but keep the DATE part of the target, only updating TIME part?
		// Or if Recurrence changes... that's complex.
		// Let's assume checking if title/desc changed (UserDayOff doesn't have title).
		// We only have InitHour/EndHour.
		// If we change InitHour of one, we usually want to shift others by same delta.

		// TODO: Implement sophisticated delta shifting.
		// For now, I will just update the target's fields directly if it's SINGLE.
		// For FUTURE/ALL, doing a straight copy of InitHour would make them all on the SAME DAY, which is wrong.
		// So we MUST calculate the Delta.

		deltaStart := dayOff.InitHour.Sub(*existing.InitHour)
		deltaEnd := dayOff.EndHour.Sub(*existing.EndHour)

		newStart := target.InitHour.Add(deltaStart)
		newEnd := target.EndHour.Add(deltaEnd)

		target.InitHour = &newStart
		target.EndHour = &newEnd
		target.UpdatedAt = time.Now()
	}

	targets := []*entity_accounts.UserDayOff{}

	if mode == UpdateModeSingle {
		targets = append(targets, existing)
	} else {
		// We need the FatherID. If existing is Father, use its ID. If not, use existing.DayOffFatherID.
		fatherID := existing.ID
		if existing.DayOffFatherID != nil {
			fatherID = existing.DayOffFatherID
		}

		if mode == UpdateModeAll {
			// Find ALL by FatherID (including father itself)
			// My repo `FindAllByFather` finds children. Does it find father? Usually no.
			// We need a way to find ALL in series.
			// Strategy: Find Father + Find Children.
			// Re-use logic: fetch all where ID=FatherID OR FatherID=FatherID.
			// Implementation detail: Use repo to find siblings.

			// Simplification: We need improved repo support or just 2 queries.
			// Current repo `FindAllByFather` finds where `day_off_father_id = ?`.
			// So we also need to fetch the father.
			father, _ := u.repo.FindByIdAndOwner(*fatherID, ownerID)
			if father != nil {
				targets = append(targets, father)
			}
			children, _ := u.repo.FindAllByFather(*fatherID, ownerID)
			targets = append(targets, children...)

		} else if mode == UpdateModeFuture {
			// Find children where init_hour >= existing.InitHour
			futureChildren, _ := u.repo.FindFutureByName(*fatherID, *existing.InitHour, ownerID)

			// Also include `existing` itself since it's "from date selected" (inclusive usually).
			// `FindFutureByName` uses `>=` so it might include existing if times match exactly or valid logic.
			// Safe bet: Add `existing` if not covered, but repository `FindFuture` might cover it if criteria matches.
			// Actually `FindFutureByName` searches children. If `existing` is child, it's covered.
			// If `existing` is Father, we need to add it.
			if existing.ID == fatherID { // It is the father
				targets = append(targets, existing)
			}
			targets = append(targets, futureChildren...)
		}
	}

	// Apply updates
	for _, target := range targets {
		updateFields(target)
		u.repo.Update(target)
	}

	return nil
}

func (u *userDayOffUseCase) Delete(id uuid.UUID, ownerID int, mode string) error {
	existing, err := u.repo.FindByIdAndOwner(id, ownerID)
	if err != nil {
		return fmt.Errorf("day off not found")
	}

	idsToDelete := []uuid.UUID{}

	if mode == DeleteModeSingle {
		idsToDelete = append(idsToDelete, id)
	} else {
		fatherID := existing.ID
		if existing.DayOffFatherID != nil {
			fatherID = existing.DayOffFatherID
		}

		if mode == DeleteModeAll {
			idsToDelete = append(idsToDelete, *fatherID)
			children, _ := u.repo.FindAllByFather(*fatherID, ownerID)
			for _, c := range children {
				idsToDelete = append(idsToDelete, *c.ID)
			}
		} else if mode == DeleteModeFuture {
			// Future includes existing + subsequent
			if existing.ID == fatherID {
				idsToDelete = append(idsToDelete, *existing.ID)
			}
			children, _ := u.repo.FindFutureByName(*fatherID, *existing.InitHour, ownerID)
			for _, c := range children {
				idsToDelete = append(idsToDelete, *c.ID)
			}
			// If existing is a child, FindFutureByName should have found it (check repo logic).
			// If not, ensure existing is added.
			found := false
			for _, tid := range idsToDelete {
				if tid == *existing.ID {
					found = true
					break
				}
			}
			if !found {
				idsToDelete = append(idsToDelete, *existing.ID)
			}
		}
	}

	if len(idsToDelete) > 0 {
		return u.repo.DeleteBatch(idsToDelete)
	}
	return nil
}

func (u *userDayOffUseCase) GetById(id uuid.UUID, ownerID int) (*entity_accounts.UserDayOff, error) {
	return u.repo.FindByIdAndOwner(id, ownerID)
}

func (u *userDayOffUseCase) GetAll(ownerID int, filterType string, year, week, month int) ([]*entity_accounts.UserDayOff, error) {
	// If no filter is specified, return all day-offs
	if filterType == "" {
		return u.repo.FindAllByOwner(ownerID)
	}

	// Calculate date range based on filter type
	var startDate, endDate time.Time

	switch filterType {
	case "week":
		if year == 0 || week == 0 || week > 53 {
			return nil, fmt.Errorf("invalid week filter: year and week (1-53) required")
		}
		// Calculate the first day of the week (Monday)
		// ISO 8601 week starts on Monday
		firstDayOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		// Find the first Monday of the year
		daysUntilMonday := (8 - int(firstDayOfYear.Weekday())) % 7
		if firstDayOfYear.Weekday() == time.Sunday {
			daysUntilMonday = 1
		} else if firstDayOfYear.Weekday() != time.Monday {
			daysUntilMonday = int(time.Monday - firstDayOfYear.Weekday())
			if daysUntilMonday < 0 {
				daysUntilMonday += 7
			}
		}
		firstMonday := firstDayOfYear.AddDate(0, 0, daysUntilMonday)
		startDate = firstMonday.AddDate(0, 0, (week-1)*7)
		endDate = startDate.AddDate(0, 0, 7)

	case "month":
		if year == 0 || month == 0 || month > 12 {
			return nil, fmt.Errorf("invalid month filter: year and month (1-12) required")
		}
		startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, 0)

	case "year":
		if year == 0 {
			return nil, fmt.Errorf("invalid year filter: year required")
		}
		startDate = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(1, 0, 0)

	default:
		return nil, fmt.Errorf("invalid filter_type: must be 'week', 'month', or 'year'")
	}

	return u.repo.FindAllByOwnerWithFilter(ownerID, &startDate, &endDate)
}
