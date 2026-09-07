export interface TeamAccountable {
  id: string;
  email: string;
}

export interface Team {
  id: string;
  name: string;
  slug: string;
  accountable: TeamAccountable;
  created_at: string;
  updated_at: string;
}

export interface CreateTeamInput {
  name: string;
  accountable_id: string;
}

export interface UpdateTeamInput {
  name?: string;
  accountable_id?: string;
}

export interface ListTeamsParams {
  name?: string;
}
