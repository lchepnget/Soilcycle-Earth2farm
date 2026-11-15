export type Envelope<T> = { success: boolean; message?: string; data: T };
export type Farmer = { id: string; full_name: string; email?: string | null; phone?: string | null; role: string; preferred_language: string };
export type Farm = { id: string; user_id: string; name: string; county: string; sub_county?: string | null; latitude: number; longitude: number; area_hectares?: number | null; water_availability?: string | null };
export type Crop = { id: string; name: string; scientific_name?: string | null; crop_family?: string | null; category?: string | null; description?: string | null };
export type SoilRecord = { id: string; farm_id: string; ph?: number | null; organic_carbon?: number | null; nitrogen?: number | null; phosphorus?: number | null; potassium?: number | null; soil_texture?: string | null; data_source: string; measurement_date?: string | null };
export type ClimateObservation = { date: string; rainfall_mm?: number; temperature_c?: number; max_temperature_c?: number; min_temperature_c?: number; humidity_percent?: number; source: string };
export type RotationResult = { id: string; plan_name: string; crop_sequence: unknown; explanation: unknown; overall_score?: number | null };
const base = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8090/api/v1').replace(/\/$/, '');
async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${base}${path}`, { ...init, headers: { 'Content-Type': 'application/json', ...(init?.headers || {}) } });
  const payload = await response.json().catch(() => ({}));
  if (!response.ok) throw new Error(payload.message || `Request failed (${response.status})`);
  return payload as T;
}
const json = (method: string, body: unknown): RequestInit => ({ method, body: JSON.stringify(body) });
export const api = {
  register: (body: { full_name: string; email?: string; phone?: string; password: string; preferred_language: string }) => request<Envelope<Farmer>>('/farmers/', json('POST', body)),
  farms: (userId: string) => request<Envelope<{ data: Farm[]; total: number }>>(`/farms/?user_id=${encodeURIComponent(userId)}&page=1&limit=100`),
  createFarm: (body: Omit<Farm, 'id'>) => request<Envelope<Farm>>('/farms/', json('POST', body)),
  crops: () => request<Envelope<{ data: Crop[]; total: number }>>('/crops/?page=1&limit=100'),
  soil: (farmId: string) => request<Envelope<{ data: SoilRecord[]; total: number }>>(`/farms/${farmId}/soil?page=1&limit=100`),
  addSoil: (body: Omit<SoilRecord, 'id'>) => request<Envelope<SoilRecord>>('/soil/', json('POST', body)),
  climate: (body: { farm_id: string; latitude: number; longitude: number; start_date: string; end_date: string }) => request<Envelope<{ source: string; observations: ClimateObservation[]; total_records: number }>>('/earthdata/climate', json('POST', body)),
  generateRotation: (body: { farm_id: string; farmer_priorities: string[]; number_of_seasons: number }) => request<Envelope<{ plan: unknown; results: RotationResult[] }>>('/rotation/', json('POST', body)),
  plans: (farmId: string) => request<Envelope<unknown[]>>(`/farms/${farmId}/rotations?limit=50&offset=0`),
};
