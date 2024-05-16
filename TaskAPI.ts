import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:3000';

interface Task {
    id?: string;
    title: string;
    description: string;
    completed: boolean;
    projectId?: string;
    assignedUserId?: string;
}

interface Project {
    id?: string;
    name: string;
    description: string;
}

class TaskService {
    static async createNewTask(taskDetails: Task): Promise<Task> {
        try {
            const response = await axios.post(`${API_BASE_URL}/tasks`, taskDetails);
            return response.data;
        } catch (error) {
            console.error('Error creating new task', error);
            throw error;
        }
    }

    static async setTaskCompleted(taskId: string): Promise<Task> {
        try {
            const response = await axios.patch(`${API_BASE_URL}/tasks/${taskId}`, { completed: true });
            return response.data;
        } catch (error) {
            console.error('Error setting task to completed', error);
            throw error;
        }
    }

    static async associateTaskWithProject(taskId: string, projectId: string): Promise<Task> {
        try {
            const response = await axios.patch(`${API_BASE_URL}/tasks/${taskId}`, { projectId });
            return response.data;
        } catch (error) {
            console.error('Error associating task with project', error);
            throw error;
        }
    }

    static async delegateTaskToUser(taskId: string, userId: string): Promise<Task> {
        try {
            const response = await axios.patch(`${API_BASE_URL}/tasks/${taskId}`, { assignedUserId: userId });
            return response.data;
        } catch (error) {
            console.error('Error delegating task to user', error);
            throw error;
        }
    }
}

class ProjectService {
    static async initiateNewProject(projectDetails: Project): Promise<Project> {
        try {
            const response = await axios.post(`${API_BASE_URL}/projects`, projectDetails);
            return response.data;
        } catch (error) {
            console.error('Error initiating new project', error);
            throw error;
        }
    }
}

export { TaskService, ProjectService };