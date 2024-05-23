import React, { useState, useEffect, ChangeEvent, FormEvent } from "react";

const API_URL = process.env.REACT_APP_API_URL;

interface Task {
  id: string;
  title: string;
  description: string;
  status: 'todo' | 'in-progress' | 'done';
}

interface Project {
  id: string;
  name: string;
  tasks: Task[];
}

const TaskBoard: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [filter, setFilter] = useState<string>('');
  const [selectedProjectId, setSelectedProjectId] = useState<string>('');

  // New task state
  const [newTask, setNewTask] = useState({
    title: '',
    description: '',
    status: 'todo' as 'todo' | 'in-progress' | 'done',
  });

  useEffect(() => {
    (async () => {
      const response = await fetch(`${API_URL}/projects`);
      const data = await response.json();
      setProjects(data);
      if(data.length > 0) setSelectedProjectId(data[0].id); // Automatically selecting the first project if available
    })();
  }, []);

  const handleFilterChange = (e: ChangeEvent<HTMLInputElement>) => {
    setFilter(e.target.value);
  };

  const handleNewTaskChange = (e: ChangeEvent<HTMLInputElement>) => {
    setNewTask({
      ...newTask,
      [e.target.name]: e.target.value,
    });
  };

  const handleStatusChange = (e: ChangeEvent<HTMLSelectElement>) => {
    setNewTask({
      ...newTask,
      status: e.target.value as 'todo' | 'in-progress' | 'done',
    });
  };

  const handleProjectChange = (e: ChangeEvent<HTMLSelectElement>) => {
    setSelectedProjectId(e.target.value);
  }

  const handleAddTask = (e: FormEvent) => {
    e.preventDefault();
    // Assuming API endpoint for adding task is /tasks with POST method
    const payload = {
      ...newTask,
      projectId: selectedProjectId,
    };

    fetch(`${API_URL}/tasks`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    }).then(async (response) => {
      const addedTask = await response.json();
      // Update corresponding project's tasks in the state
      setProjects(prevProjects =>
        prevProjects.map(project => {
          if (project.id === selectedProjectId) {
            return {
              ...project,
              tasks: [...project.tasks, addedTask],
            };
          }
          return project;
        }),
      );
    });

    // Resetting form
    setNewTask({ title: '', description: '', status: 'todo' });
  };

  const filteredProjects = projects
    .map(project => ({
      ...project,
      tasks: project.tasks.filter(task => task.title.toLowerCase().includes(filter.toLowerCase()) || task.description.toLowerCase().includes(filter.toLowerCase())),
    }))
    .filter(project => project.tasks.length > 0);

  return (
    <div>
      <input type="text" placeholder="Filter tasks" value={filter} onChange={handleFilterChange} />
      <form onSubmit={handleAddTask}>
        <h2>Add New Task</h2>
        <select value={selectedProjectId} onChange={handleProjectChange} required>
          {projects.map((project) => (
            <option key={project.id} value={project.id}>{project.name}</option>
          ))}
        </select>
        <input name="title" value={newTask.title} onChange={handleNewTaskChange} placeholder="Title" required/>
        <input name="description" value={newTask.description} onChange={handleNewTaskChange} placeholder="Description" required/>
        <select value={newTask.status} onChange={handleStatusChange}>
          <option value="todo">To Do</option>
          <option value="in-progress">In Progress</option>
          <option value="done">Done</option>
        </select>
        <button type="submit">Add Task</button>
      </form>
      {filteredProjects.map(project => (
        <div key={project.id}>
          <h2>{project.name}</h2>
          {project.tasks.map(task => (
            <div key={task.id}>
              <h3>{task.title}</h3>
              <p>{task.description}</p>
              <p>Status: {task.status}</p>
            </div>
          ))}
        </div>
      ))}
    </div>
  );
};

export default TaskBoard;