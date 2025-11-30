"use server";

import { getMyGroups as fetchMyGroups, getGroupById as fetchGroupById, getMockGroups as fetchMockGroups } from "../../mock-data/group";
import { getMockPosts } from "../../mock-data/posts";
import { getMockEvents } from "../../mock-data/events";
import { getUserByID } from "../../mock-data/users";

export async function getMyGroups() {
    // simulate api delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return fetchMyGroups();
}

export async function getGroupById(id) {
    // simulate api delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return fetchGroupById(id);
}

export async function getAllGroups() {
    // simulate api delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return fetchMockGroups();
}

export async function getGroupPosts() {
    // simulate api delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    // For mock purposes, return all posts
    // In a real app, this would filter by groupId
    return getMockPosts();
}

export async function getGroupMembers() {
    // simulate api delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    // For mock purposes, return a few users
    const userIds = ["1", "2", "3", "4", "5"];
    return userIds.map(id => getUserByID(id));
}

export async function getGroupEvents() {
    // simulate api delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return getMockEvents();
}
