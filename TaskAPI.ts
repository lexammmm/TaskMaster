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

class TaskAPI {
    static async createTask(taskDetails: Task): Promise<Task> {
        try {
            const response = await axios.post(`${API_BASE_URL}/tasks`, taskDetails);
            return response.data;
        } catch (error) {
            console.error('Error creating task', error);
            throw error;
        }
    }

    static async markTaskAsCompleted(taskId: string): Promise<Task> {
        try {
            const response = await axios.patch(`${API_BASE_URL}/tasks/${taskId}`, { completed: true });
            return response.data;
        } catch (error) {
            console.error('Error marking task as completed', error);
            throw error;
        }
    }

    static async linkTaskToProject(taskId: string, projectId: string): Promise<Task> {
        try {
            const response = await axios.patch(`${API_BASE_URL}/tasks/${taskId}`, { projectId: projectId });
            return response.data;
        } catch (error) {
            console.error('Error linking task to project', error);
            throw error;
        }
    }

    static async assignTaskToUser(taskId: string, userId: string): Promise<Task> {
        try {
            const response = await axios.patch(`${API_BASE_URL}/tasks/${taskId}`, { assignedUserId: userId });
            return response.data;
        } catch (error) {
            console.error('Error assigning task to user', error);
            throw error;
        }
    }
}

class ProjectAPI {
    static async createProject(projectDetails: Project): Promise<Project> {
        try {
            const response = await axios.post(`${API_BASE_URL}/projects`, projectDetails);
            return response.data;
        } catch (error) {
            console.error('Error creating project', error);
            throw error;
        }
    }
}

export { TaskAPI as TaskService, ProjectAPI as ProjectService };