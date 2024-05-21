import React, { useState, useEffect, ChangeEvent } from "react";

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

  useEffect(() => {
    (async () => {
      const response = await fetch(`${API_URL}/projects`);
      const data = await response.json();
      setProjects(data);
    })();
  }, []);

  const handleFilterChange = (e: ChangeEvent<HTMLInputElement>) => {
    setFilter(e.target.value);
  };

  const filteredProjects = projects
    .map(project => ({
      ...project,
      tasks: project.tasks.filter(task => task.title.includes(filter) || task.description.includes(filter)),
    }))
    .filter(project => project.tasks.length > 0);

  return (
    <div>
      <input type="text" placeholder="Filter tasks" value={filter} onChange={handleFilterChange} />
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