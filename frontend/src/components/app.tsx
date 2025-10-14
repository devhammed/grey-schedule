import React, {useCallback, useMemo, useState} from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createAppointment, deleteAppointment, listAppointments } from '../utils/api.ts'
import type { CreateAppointmentInput } from '../types.ts'
import { fromLocalInputValue, toLocalInputValue } from '../utils/date.ts'

export default function App() {
    const qc = useQueryClient()

    const { data: appointments, isLoading, error } = useQuery({
        queryKey: ['appointments'],
        queryFn: listAppointments,
    })

    const [title, setTitle] = useState('')
    const [startLocal, setStartLocal] = useState(() => toLocalInputValue(new Date()))
    const [endLocal, setEndLocal] = useState(() => toLocalInputValue(new Date(Date.now() + 30 * 60 * 1000)))

    const createMut = useMutation({
        mutationFn: (input: CreateAppointmentInput) => createAppointment(input),
        onSuccess: () => {
            void qc.invalidateQueries({ queryKey: ['appointments'] })
            setTitle('')
        },
    })

    const deleteMut = useMutation({
        mutationFn: (id: string) => deleteAppointment(id),
        onSuccess: () => qc.invalidateQueries({ queryKey: ['appointments'] }),
    })

    const sorted = useMemo(() => {
        return (appointments ?? []).slice().sort((a, b) => a.start.localeCompare(b.start))
    }, [appointments])

    const onSubmit = useCallback((e: React.FormEvent) => {
        e.preventDefault()

        const start = fromLocalInputValue(startLocal)

        const end = fromLocalInputValue(endLocal)

        createMut.mutate({ title: title.trim(), start, end })
    }, [startLocal, endLocal, title, createMut]);

    return (
        <div className="max-w-3xl mx-auto mt-8 p-4 font-sans">
            <h1 className="text-3xl font-bold mb-6 text-center">Schedule Management</h1>

            <section className="mb-8 p-6 border border-gray-300 rounded-xl bg-white shadow-sm">
                <h2 className="text-xl font-semibold mb-4">Create Appointment</h2>

                <form onSubmit={onSubmit} className="grid gap-4">
                    <label className="flex flex-col">
                        <span className="mb-1 font-medium">Title</span>
                        <input
                            value={title}
                            onChange={e => setTitle(e.target.value)}
                            placeholder="Team sync"
                            required
                            className="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring focus:ring-primary/80"
                        />
                    </label>

                    <label className="flex flex-col">
                        <span className="mb-1 font-medium">Start</span>
                        <input
                            type="datetime-local"
                            value={startLocal}
                            onChange={e => setStartLocal(e.target.value)}
                            required
                            className="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring focus:ring-primary/80"
                        />
                    </label>

                    <label className="flex flex-col">
                        <span className="mb-1 font-medium">End</span>
                        <input
                            type="datetime-local"
                            value={endLocal}
                            onChange={e => setEndLocal(e.target.value)}
                            required
                            className="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring focus:ring-primary/80"
                        />
                    </label>

                    <button
                        type="submit"
                        disabled={createMut.isPending}
                        className="bg-primary text-white rounded-md py-2 font-medium hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {createMut.isPending ? 'Creating...' : 'Create'}
                    </button>

                    {createMut.isError && (
                        <div className="text-red-600 text-sm">
                            {(createMut.error as any)?.response?.data?.error || 'An error occurred'}
                        </div>
                    )}

                    {createMut.isSuccess && <div className="text-green-600 text-sm">Created!</div>}
                </form>
            </section>

            <section className="p-6 border border-gray-300 rounded-xl bg-white shadow-sm">
                <h2 className="text-xl font-semibold mb-4">Appointments</h2>

                {isLoading && <p>Loading…</p>}

                {error && <p className="text-red-600">Failed to load</p>}

                {!isLoading && sorted.length === 0 && <p>No appointments yet</p>}

                {sorted.length > 0 && (
                    <ul className="space-y-3">
                        {sorted.map((a) => (
                            <li
                                key={a.id}
                                className="flex justify-between items-center p-3 border border-gray-200 rounded-md hover:bg-gray-50"
                            >
                                <div>
                                    <div className="font-semibold">{a.title}</div>
                                    <div className="text-gray-600 text-sm">
                                        {new Date(a.start).toLocaleString()} → {new Date(a.end).toLocaleString()}
                                    </div>
                                </div>
                                <button
                                    onClick={() => deleteMut.mutate(a.id)}
                                    disabled={deleteMut.isPending}
                                    className="text-sm bg-red-500 text-white px-3 py-1 rounded-md hover:bg-red-600 disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    Delete
                                </button>
                            </li>
                        ))}
                    </ul>
                )}
            </section>
        </div>
    )
}