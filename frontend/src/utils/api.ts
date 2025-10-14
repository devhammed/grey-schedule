import axios from 'axios'
import type { Appointment,CreateAppointmentInput } from '../types'

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

export async function listAppointments(): Promise<Appointment[]> {
  const res = await axios.get(`${API_BASE}/api/appointments`)
  return res.data
}

export async function createAppointment(input: CreateAppointmentInput): Promise<Appointment> {
  const res = await axios.post(`${API_BASE}/api/appointments`, input)
  return res.data
}

export async function deleteAppointment(id: string): Promise<void> {
  await axios.delete(`${API_BASE}/api/appointments/${id}`)
}
