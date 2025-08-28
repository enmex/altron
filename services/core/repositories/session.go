package repositories

import (
	common "altron/common/models"
	"altron/core/interfaces"
	"altron/core/models"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/packet"
	"altron/core/repositories/ent/service"
	"altron/core/repositories/ent/session"
	"altron/core/repositories/ent/sessionfilter"
	"altron/core/repositories/ent/user"
	"altron/core/repositories/ent/workspace"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ interfaces.SessionRepository = (*SessionRepository)(nil)

type SessionRepository struct {
	Client *ent.Client
}

func NewSessionRepository(client *ent.Client) *SessionRepository {
	return &SessionRepository{
		Client: client,
	}
}

func (r *SessionRepository) CreateSessions(ctx context.Context, workspaceID uuid.UUID, sessions []*common.Session) error {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		sessionEnt, err := tx.Session.
			Create().
			SetID(session.ID).
			SetWorkspaceID(workspaceID).
			SetClientHost(session.ClientHost).
			SetServerPort(uint32(session.ServerPort)).
			SetProtocol(session.Protocol).
			SetSentAt(session.SentAt).
			SetTTL(session.TTL).
			SetPacketsCount(session.PacketsCount).
			SetNillableClientUserAgent(session.ClientUserAgent).
			SetAverageResponseTime(session.AverageResponseTime).
			SetRequestsNumber(session.RequestsNumber).
			Save(ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
		for _, f := range session.MatchedFilters {
			sessionFilterEnt, err := tx.SessionFilter.
				Create().
				SetFilterID(f.ID).
				SetSessionID(sessionEnt.ID).
				SetMatchesCount(f.MatchesCount).
				Save(ctx)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
			for _, packetIdx := range f.MatchedPackets {
				err = tx.MatchedPacket.Create().
					SetSessionFilter(sessionFilterEnt).
					SetPacketIdx(packetIdx).
					Exec(ctx)
				if err != nil {
					if err := tx.Rollback(); err != nil {
						return err
					}
					return err
				}
			}
		}

		for _, packet := range session.Packets {
			_, err = tx.Packet.
				Create().
				SetSentAt(packet.SentAt).
				SetSession(sessionEnt).
				SetIsRequest(packet.IsRequest).
				SetPayload(packet.Payload).
				Save(ctx)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *SessionRepository) GetSessionsByWorkspace(ctx context.Context, workspaceID uuid.UUID, filterID *uuid.UUID, paginationIndex int) ([]*ent.Session, error) {
	predicate := session.HasWorkspaceWith(workspace.IDEQ(workspaceID))
	if filterID != nil {
		predicate = session.And(predicate, session.HasSessionFiltersWith(sessionfilter.FilterIDEQ(*filterID)))
	}
	sessions, err := r.Client.Session.
		Query().
		WithSessionFilters(func(sfq *ent.SessionFilterQuery) {
			sfq.WithFilter().WithMatchedPackets()
		}).
		Offset(paginationIndex * 100).
		Order(ent.Asc("sent_at")).
		Limit(100).
		Where(predicate).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) GetSessionPackets(ctx context.Context, sessionID uuid.UUID) ([]*ent.Packet, error) {
	return r.Client.Packet.
		Query().
		Where(packet.HasSessionWith(session.IDEQ(sessionID))).
		All(ctx)
}

func (r *SessionRepository) DeleteAllSessions(ctx context.Context, workspaceID uuid.UUID) error {
	_, err := r.Client.Session.
		Delete().
		Where(session.HasWorkspaceWith(workspace.IDEQ(workspaceID))).
		Exec(ctx)
	return err
}

func (r *SessionRepository) CountWorkspaceSessions(ctx context.Context, workspaceID uuid.UUID, filterID *uuid.UUID) (int, error) {
	predicate := session.HasWorkspaceWith(workspace.IDEQ(workspaceID))
	if filterID != nil {
		predicate = session.And(predicate, session.HasSessionFiltersWith(sessionfilter.FilterIDEQ(*filterID)))
	}
	return r.Client.Session.
		Query().
		Where(predicate).
		Count(ctx)
}

func (r *SessionRepository) GetSession(ctx context.Context, sessionID uuid.UUID) (*ent.Session, error) {
	return r.Client.Session.
		Query().
		WithPackets().
		WithSessionFilters(func(sfq *ent.SessionFilterQuery) {
			sfq.WithFilter().WithMatchedPackets()
		}).
		Where(session.IDEQ(sessionID)).
		Only(ctx)
}

func (r *SessionRepository) CreateSession(ctx context.Context, userID uuid.UUID, session *common.Session) error {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return err
	}
	userWorkspaceEnt, err := tx.Workspace.
		Query().
		Where(workspace.NameEQ(userID.String())).
		Only(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	sessionEnt, err := tx.Session.
		Create().
		SetID(session.ID).
		SetClientHost(session.ClientHost).
		SetWorkspace(userWorkspaceEnt).
		SetServerPort(uint32(session.ServerPort)).
		SetProtocol(session.Protocol).
		SetSentAt(session.SentAt).
		SetTTL(session.TTL).
		SetPacketsCount(session.PacketsCount).
		SetNillableClientUserAgent(session.ClientUserAgent).
		SetAverageResponseTime(session.AverageResponseTime).
		SetRequestsNumber(session.RequestsNumber).
		Save(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		if strings.Contains(err.Error(), "duplicate") {
			return models.ErrorDuplicateSession
		}
		return err
	}
	for _, f := range session.MatchedFilters {
		sessionFilterEnt, err := tx.SessionFilter.
			Create().
			SetFilterID(f.ID).
			SetSessionID(sessionEnt.ID).
			SetMatchesCount(f.MatchesCount).
			Save(ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
		for _, packetIdx := range f.MatchedPackets {
			err = tx.MatchedPacket.Create().
				SetSessionFilter(sessionFilterEnt).
				SetPacketIdx(packetIdx).
				Exec(ctx)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
		}
	}

	for _, packet := range session.Packets {
		_, err = tx.Packet.
			Create().
			SetSentAt(packet.SentAt).
			SetSession(sessionEnt).
			SetIsRequest(packet.IsRequest).
			SetPayload(packet.Payload).
			Save(ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return r.Client.Session.
		DeleteOneID(sessionID).
		Exec(ctx)
}

func (r *SessionRepository) GetSessionsByIDs(ctx context.Context, sessionIDs []uuid.UUID) ([]*ent.Session, error) {
	sessionsEnt, err := r.Client.Session.Query().
		WithSessionFilters(func(sfq *ent.SessionFilterQuery) {
			sfq.WithFilter().WithMatchedPackets()
		}).
		WithPackets().
		Order(ent.Asc("sent_at")).
		Where(session.IDIn(sessionIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return sessionsEnt, nil
}

func (r *SessionRepository) GetPcapSessions(ctx context.Context, pcapWorkspaceID uuid.UUID, paginationIndex int) ([]*ent.Session, error) {
	sessions, err := r.Client.Session.
		Query().
		WithSessionFilters(func(sfq *ent.SessionFilterQuery) {
			sfq.WithFilter()
		}).
		Offset(paginationIndex * 100).
		Order(ent.Asc("sent_at")).
		Limit(100).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) CreatePcapSession(ctx context.Context, pcapWorkspaceID uuid.UUID, session *common.Session) error {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return err
	}
	sessionEnt, err := tx.Session.
		Create().
		SetID(session.ID).
		AddPcapWorkspaceIDs(pcapWorkspaceID).
		SetClientHost(session.ClientHost).
		SetServerPort(uint32(session.ServerPort)).
		SetProtocol(session.Protocol).
		SetSentAt(session.SentAt).
		SetTTL(session.TTL).
		SetPacketsCount(session.PacketsCount).
		SetNillableClientUserAgent(session.ClientUserAgent).
		SetAverageResponseTime(session.AverageResponseTime).
		SetRequestsNumber(session.RequestsNumber).
		Save(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	for _, f := range session.MatchedFilters {
		err := tx.SessionFilter.
			Create().
			SetFilterID(f.ID).
			SetSessionID(sessionEnt.ID).
			SetMatchesCount(f.MatchesCount).
			Exec(ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	for _, packet := range session.Packets {
		_, err = tx.Packet.
			Create().
			SetSentAt(packet.SentAt).
			SetSession(sessionEnt).
			SetIsRequest(packet.IsRequest).
			SetPayload(packet.Payload).
			Save(ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *SessionRepository) CountSessions(ctx context.Context, userID uuid.UUID) (int, error) {
	return r.Client.Session.
		Query().
		Where(session.HasWorkspaceWith(workspace.Or(
			workspace.HasServiceWith(service.HasUserWith(user.IDEQ(userID))),
			workspace.NameEQ(userID.String()),
		))).
		Count(ctx)
}

func (r *SessionRepository) GetSessions(ctx context.Context, userID uuid.UUID, paginationIndex int) ([]*ent.Session, error) {
	return r.Client.Session.
		Query().
		WithPackets(func(pq *ent.PacketQuery) {
			pq.Order(ent.Asc("sent_at"))
		}).
		Where(session.HasWorkspaceWith(workspace.Or(
			workspace.HasServiceWith(service.HasUserWith(user.IDEQ(userID))),
			workspace.NameEQ(userID.String()),
		))).
		Offset(paginationIndex * 100).
		Limit(100).
		All(ctx)
}