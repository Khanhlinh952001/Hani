export type UserGender = "male" | "female" | "other";

export const GENDER_OPTIONS: {
  id: UserGender;
  label: string;
  desc: string;
}[] = [
  { id: "male", label: "Nam", desc: "Hani gọi bạn kiểu 오빠" },
  { id: "female", label: "Nữ", desc: "Hani gọi bạn kiểu 언니/누나" },
  { id: "other", label: "Khác", desc: "Hani dùng tên bạn, trung tính" },
];
