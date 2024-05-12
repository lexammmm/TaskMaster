import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:3000';

interface Task {
    id?: string;
    title: string;
    description: string;
    completed: boolean;
    projectId?: string;
    assignedTo?: string;
}

interface Project {
    id?: string;
    name: string;
    description: string;
}

class TaskService {
    static async addTask(task: Task): Promise<Task> {
        try {
            const response = await axios.post(`${API_BASE_URL}/tasks`, task);
            return response.data;
        } catch (error) {
            console.error('Error adding task', error);
            throw error;
        }
    }

    static async completeTask(taskId: string): Promise<Task> {
        try {
            const response = await axios.patch(`${API_BASE_URL}/tasks/${taskId}`, { completed: true });
            return response.data;
        } catch (error) {
            console.error('Error completing task', error);
            throw error;
        }
    }

    static async assignTaskToProject(taskId: string, projectId: string): Promise<Task> {
        try {
            const response = await axios.patch(`${API_BASE_URL}/tasks/${taskId}`, { projectId: projectId });
            return response.data;
        } catch (error) {
            console.error('Error assigning task to project', error);
            throw error;
        }
    }

    static async assignTaskToUser(taskId: string, userId: string): Promise<Task> {
        try {
            const response = await axios.patch(`${API_BASE_URL}/tasks/${taskId}`, { assignedTo: userId });
            return response.data;
        } catch (error) {
            console.error('Error assigning task to user', error);
            throw error;
        }
    }
}

class ProjectService {
    static async createProject(project: Project): Promise<Project> {
        try {
            const response = await axios.post(`${API_BASE_URL}/projects`, project);
            return response.data;
        } catch (error) {
            console.error('Error creating project', error);
            throw error;
        }
    }
}

export { TaskService, ProjectService };