export type Appointment = {
  id: string
  title: string
  start: string
  end: string
  createdAt: string
}

export type CreateAppointmentInput = {
    title: string
    start: string
    end: string
}
