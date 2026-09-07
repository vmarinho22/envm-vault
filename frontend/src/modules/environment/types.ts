export interface Environment {
  id: string;
  title: string;
  created_at: string;
  updated_at: string;
}

export interface CreateEnvironmentInput {
  title: string;
}

export interface UpdateEnvironmentInput {
  title?: string;
}

export interface ListEnvironmentsParams {
  title?: string;
}
