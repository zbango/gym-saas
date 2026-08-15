import React from "react";
import { SafeAreaView, Text, View } from "react-native";
import { appName, desktopVersion } from "@gym-saas/shared";

export default function App() {
  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: "#f7efe4" }}>
      <View style={{ padding: 24, gap: 12 }}>
        <Text style={{ fontSize: 28, fontWeight: "700", color: "#34220b" }}>{appName}</Text>
        <Text style={{ fontSize: 18, color: "#6f5634" }}>Expo placeholder host</Text>
        <Text style={{ color: "#5d4a31" }}>Shared version: {desktopVersion}</Text>
      </View>
    </SafeAreaView>
  );
}
