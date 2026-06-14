export interface Team {
  id: string;
  name: string;
  code: string;
  flag_url: string;
  group_code: string;
  group_points: number;
  goals_for: number;
  goals_against: number;
  goal_difference: number;
  group_rank: number;
}

export interface Match {
  id: string;
  phase: 'group' | 'r32' | 'r16' | 'qf' | 'sf' | 'final' | 'third';
  group_code?: string;
  match_number: number;
  home_team: string;
  away_team: string;
  kickoff: string;
  venue: string;
  city: string;
  score_home?: number;
  score_away?: number;
  winner?: string;
  status: 'upcoming' | 'live' | 'finished';
  match_day?: number;
  expand?: {
    home_team?: Team;
    away_team?: Team;
    winner?: Team;
  };
}

export interface Prediction {
  id: string;
  user: string;
  match: string;
  prediction_group: string;
  predicted_home: number;
  predicted_away: number;
  points_awarded: number;
  is_locked: boolean;
  submitted_at: string;
}

export interface PredictionGroup {
  id: string;
  name: string;
  owner: string;
  invite_code?: string;
  invite_expires_at?: string;
  is_public: boolean;
  expand?: {
    owner?: {
      id: string;
      username: string;
    };
  };
}

export interface LeaderboardUser {
  userId: string;
  username: string;
  avatarUrl: string;
  avatar: string;
  totalPoints: number;
  rank: number;
}
