// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common/pb_msg/ai_camera.proto

#ifndef GOOGLE_PROTOBUF_INCLUDED_common_2fpb_5fmsg_2fai_5fcamera_2eproto_2epb_2eh
#define GOOGLE_PROTOBUF_INCLUDED_common_2fpb_5fmsg_2fai_5fcamera_2eproto_2epb_2eh

#include <limits>
#include <string>
#include <type_traits>

#include "google/protobuf/port_def.inc"
#if PROTOBUF_VERSION < 4024000
#error "This file was generated by a newer version of protoc which is"
#error "incompatible with your Protocol Buffer headers. Please update"
#error "your headers."
#endif  // PROTOBUF_VERSION

#if 4024000 < PROTOBUF_MIN_PROTOC_VERSION
#error "This file was generated by an older version of protoc which is"
#error "incompatible with your Protocol Buffer headers. Please"
#error "regenerate this file with a newer version of protoc."
#endif  // PROTOBUF_MIN_PROTOC_VERSION
#include "google/protobuf/port_undef.inc"
#include "google/protobuf/io/coded_stream.h"
#include "google/protobuf/arena.h"
#include "google/protobuf/arenastring.h"
#include "google/protobuf/generated_message_tctable_decl.h"
#include "google/protobuf/generated_message_util.h"
#include "google/protobuf/metadata_lite.h"
#include "google/protobuf/generated_message_reflection.h"
#include "google/protobuf/message.h"
#include "google/protobuf/repeated_field.h"  // IWYU pragma: export
#include "google/protobuf/extension_set.h"  // IWYU pragma: export
#include "google/protobuf/unknown_field_set.h"
#include "common/pb_msg/common.pb.h"
// @@protoc_insertion_point(includes)

// Must be included last.
#include "google/protobuf/port_def.inc"

#define PROTOBUF_INTERNAL_EXPORT_common_2fpb_5fmsg_2fai_5fcamera_2eproto

namespace google {
namespace protobuf {
namespace internal {
class AnyMetadata;
}  // namespace internal
}  // namespace protobuf
}  // namespace google

// Internal implementation detail -- do not use these members.
struct TableStruct_common_2fpb_5fmsg_2fai_5fcamera_2eproto {
  static const ::uint32_t offsets[];
};
extern const ::google::protobuf::internal::DescriptorTable
    descriptor_table_common_2fpb_5fmsg_2fai_5fcamera_2eproto;
namespace pb_msg {
class AiCamera;
struct AiCameraDefaultTypeInternal;
extern AiCameraDefaultTypeInternal _AiCamera_default_instance_;
}  // namespace pb_msg
namespace google {
namespace protobuf {
}  // namespace protobuf
}  // namespace google

namespace pb_msg {

// ===================================================================


// -------------------------------------------------------------------

class AiCamera final :
    public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:pb_msg.AiCamera) */ {
 public:
  inline AiCamera() : AiCamera(nullptr) {}
  ~AiCamera() override;
  template<typename = void>
  explicit PROTOBUF_CONSTEXPR AiCamera(::google::protobuf::internal::ConstantInitialized);

  AiCamera(const AiCamera& from);
  AiCamera(AiCamera&& from) noexcept
    : AiCamera() {
    *this = ::std::move(from);
  }

  inline AiCamera& operator=(const AiCamera& from) {
    CopyFrom(from);
    return *this;
  }
  inline AiCamera& operator=(AiCamera&& from) noexcept {
    if (this == &from) return *this;
    if (GetOwningArena() == from.GetOwningArena()
  #ifdef PROTOBUF_FORCE_COPY_IN_MOVE
        && GetOwningArena() != nullptr
  #endif  // !PROTOBUF_FORCE_COPY_IN_MOVE
    ) {
      InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }

  inline const ::google::protobuf::UnknownFieldSet& unknown_fields() const {
    return _internal_metadata_.unknown_fields<::google::protobuf::UnknownFieldSet>(::google::protobuf::UnknownFieldSet::default_instance);
  }
  inline ::google::protobuf::UnknownFieldSet* mutable_unknown_fields() {
    return _internal_metadata_.mutable_unknown_fields<::google::protobuf::UnknownFieldSet>();
  }

  static const ::google::protobuf::Descriptor* descriptor() {
    return GetDescriptor();
  }
  static const ::google::protobuf::Descriptor* GetDescriptor() {
    return default_instance().GetMetadata().descriptor;
  }
  static const ::google::protobuf::Reflection* GetReflection() {
    return default_instance().GetMetadata().reflection;
  }
  static const AiCamera& default_instance() {
    return *internal_default_instance();
  }
  static inline const AiCamera* internal_default_instance() {
    return reinterpret_cast<const AiCamera*>(
               &_AiCamera_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    0;

  friend void swap(AiCamera& a, AiCamera& b) {
    a.Swap(&b);
  }
  inline void Swap(AiCamera* other) {
    if (other == this) return;
  #ifdef PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() != nullptr &&
        GetOwningArena() == other->GetOwningArena()) {
   #else  // PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() == other->GetOwningArena()) {
  #endif  // !PROTOBUF_FORCE_COPY_IN_SWAP
      InternalSwap(other);
    } else {
      ::google::protobuf::internal::GenericSwap(this, other);
    }
  }
  void UnsafeArenaSwap(AiCamera* other) {
    if (other == this) return;
    ABSL_DCHECK(GetOwningArena() == other->GetOwningArena());
    InternalSwap(other);
  }

  // implements Message ----------------------------------------------

  AiCamera* New(::google::protobuf::Arena* arena = nullptr) const final {
    return CreateMaybeMessage<AiCamera>(arena);
  }
  using ::google::protobuf::Message::CopyFrom;
  void CopyFrom(const AiCamera& from);
  using ::google::protobuf::Message::MergeFrom;
  void MergeFrom( const AiCamera& from) {
    AiCamera::MergeImpl(*this, from);
  }
  private:
  static void MergeImpl(::google::protobuf::Message& to_msg, const ::google::protobuf::Message& from_msg);
  public:
  PROTOBUF_ATTRIBUTE_REINITIALIZES void Clear() final;
  bool IsInitialized() const final;

  ::size_t ByteSizeLong() const final;
  const char* _InternalParse(const char* ptr, ::google::protobuf::internal::ParseContext* ctx) final;
  ::uint8_t* _InternalSerialize(
      ::uint8_t* target, ::google::protobuf::io::EpsCopyOutputStream* stream) const final;
  int GetCachedSize() const final { return _impl_._cached_size_.Get(); }

  private:
  void SharedCtor(::google::protobuf::Arena* arena);
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(AiCamera* other);

  private:
  friend class ::google::protobuf::internal::AnyMetadata;
  static ::absl::string_view FullMessageName() {
    return "pb_msg.AiCamera";
  }
  protected:
  explicit AiCamera(::google::protobuf::Arena* arena);
  public:

  static const ClassData _class_data_;
  const ::google::protobuf::Message::ClassData*GetClassData() const final;

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  enum : int {
    kNameFieldNumber = 1,
    kPositionFieldNumber = 2,
    kUpFieldNumber = 3,
    kLookAtFieldNumber = 4,
    kHorizontalFOVFieldNumber = 5,
    kClipPlaneNearFieldNumber = 6,
    kClipPlaneFarFieldNumber = 7,
    kAspectFieldNumber = 8,
    kOrthographicWidthFieldNumber = 9,
  };
  // string Name = 1;
  void clear_name() ;
  const std::string& name() const;
  template <typename Arg_ = const std::string&, typename... Args_>
  void set_name(Arg_&& arg, Args_... args);
  std::string* mutable_name();
  PROTOBUF_NODISCARD std::string* release_name();
  void set_allocated_name(std::string* ptr);

  private:
  const std::string& _internal_name() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_name(
      const std::string& value);
  std::string* _internal_mutable_name();

  public:
  // .pb_msg.AiVector3D Position = 2;
  bool has_position() const;
  void clear_position() ;
  const ::pb_msg::AiVector3D& position() const;
  PROTOBUF_NODISCARD ::pb_msg::AiVector3D* release_position();
  ::pb_msg::AiVector3D* mutable_position();
  void set_allocated_position(::pb_msg::AiVector3D* value);
  void unsafe_arena_set_allocated_position(::pb_msg::AiVector3D* value);
  ::pb_msg::AiVector3D* unsafe_arena_release_position();

  private:
  const ::pb_msg::AiVector3D& _internal_position() const;
  ::pb_msg::AiVector3D* _internal_mutable_position();

  public:
  // .pb_msg.AiVector3D Up = 3;
  bool has_up() const;
  void clear_up() ;
  const ::pb_msg::AiVector3D& up() const;
  PROTOBUF_NODISCARD ::pb_msg::AiVector3D* release_up();
  ::pb_msg::AiVector3D* mutable_up();
  void set_allocated_up(::pb_msg::AiVector3D* value);
  void unsafe_arena_set_allocated_up(::pb_msg::AiVector3D* value);
  ::pb_msg::AiVector3D* unsafe_arena_release_up();

  private:
  const ::pb_msg::AiVector3D& _internal_up() const;
  ::pb_msg::AiVector3D* _internal_mutable_up();

  public:
  // .pb_msg.AiVector3D LookAt = 4;
  bool has_lookat() const;
  void clear_lookat() ;
  const ::pb_msg::AiVector3D& lookat() const;
  PROTOBUF_NODISCARD ::pb_msg::AiVector3D* release_lookat();
  ::pb_msg::AiVector3D* mutable_lookat();
  void set_allocated_lookat(::pb_msg::AiVector3D* value);
  void unsafe_arena_set_allocated_lookat(::pb_msg::AiVector3D* value);
  ::pb_msg::AiVector3D* unsafe_arena_release_lookat();

  private:
  const ::pb_msg::AiVector3D& _internal_lookat() const;
  ::pb_msg::AiVector3D* _internal_mutable_lookat();

  public:
  // float HorizontalFOV = 5;
  void clear_horizontalfov() ;
  float horizontalfov() const;
  void set_horizontalfov(float value);

  private:
  float _internal_horizontalfov() const;
  void _internal_set_horizontalfov(float value);

  public:
  // float ClipPlaneNear = 6;
  void clear_clipplanenear() ;
  float clipplanenear() const;
  void set_clipplanenear(float value);

  private:
  float _internal_clipplanenear() const;
  void _internal_set_clipplanenear(float value);

  public:
  // float ClipPlaneFar = 7;
  void clear_clipplanefar() ;
  float clipplanefar() const;
  void set_clipplanefar(float value);

  private:
  float _internal_clipplanefar() const;
  void _internal_set_clipplanefar(float value);

  public:
  // float Aspect = 8;
  void clear_aspect() ;
  float aspect() const;
  void set_aspect(float value);

  private:
  float _internal_aspect() const;
  void _internal_set_aspect(float value);

  public:
  // float OrthographicWidth = 9;
  void clear_orthographicwidth() ;
  float orthographicwidth() const;
  void set_orthographicwidth(float value);

  private:
  float _internal_orthographicwidth() const;
  void _internal_set_orthographicwidth(float value);

  public:
  // @@protoc_insertion_point(class_scope:pb_msg.AiCamera)
 private:
  class _Internal;

  friend class ::google::protobuf::internal::TcParser;
  static const ::google::protobuf::internal::TcParseTable<4, 9, 3, 36, 2> _table_;
  template <typename T> friend class ::google::protobuf::Arena::InternalHelper;
  typedef void InternalArenaConstructable_;
  typedef void DestructorSkippable_;
  struct Impl_ {
    ::google::protobuf::internal::HasBits<1> _has_bits_;
    mutable ::google::protobuf::internal::CachedSize _cached_size_;
    ::google::protobuf::internal::ArenaStringPtr name_;
    ::pb_msg::AiVector3D* position_;
    ::pb_msg::AiVector3D* up_;
    ::pb_msg::AiVector3D* lookat_;
    float horizontalfov_;
    float clipplanenear_;
    float clipplanefar_;
    float aspect_;
    float orthographicwidth_;
    PROTOBUF_TSAN_DECLARE_MEMBER;
  };
  union { Impl_ _impl_; };
  friend struct ::TableStruct_common_2fpb_5fmsg_2fai_5fcamera_2eproto;
};

// ===================================================================




// ===================================================================


#ifdef __GNUC__
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wstrict-aliasing"
#endif  // __GNUC__
// -------------------------------------------------------------------

// AiCamera

// string Name = 1;
inline void AiCamera::clear_name() {
  _impl_.name_.ClearToEmpty();
}
inline const std::string& AiCamera::name() const {
  // @@protoc_insertion_point(field_get:pb_msg.AiCamera.Name)
  return _internal_name();
}
template <typename Arg_, typename... Args_>
inline PROTOBUF_ALWAYS_INLINE void AiCamera::set_name(Arg_&& arg,
                                                     Args_... args) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.name_.Set(static_cast<Arg_&&>(arg), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:pb_msg.AiCamera.Name)
}
inline std::string* AiCamera::mutable_name() {
  std::string* _s = _internal_mutable_name();
  // @@protoc_insertion_point(field_mutable:pb_msg.AiCamera.Name)
  return _s;
}
inline const std::string& AiCamera::_internal_name() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.name_.Get();
}
inline void AiCamera::_internal_set_name(const std::string& value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.name_.Set(value, GetArenaForAllocation());
}
inline std::string* AiCamera::_internal_mutable_name() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  return _impl_.name_.Mutable( GetArenaForAllocation());
}
inline std::string* AiCamera::release_name() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  // @@protoc_insertion_point(field_release:pb_msg.AiCamera.Name)
  return _impl_.name_.Release();
}
inline void AiCamera::set_allocated_name(std::string* value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  _impl_.name_.SetAllocated(value, GetArenaForAllocation());
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        if (_impl_.name_.IsDefault()) {
          _impl_.name_.Set("", GetArenaForAllocation());
        }
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:pb_msg.AiCamera.Name)
}

// .pb_msg.AiVector3D Position = 2;
inline bool AiCamera::has_position() const {
  bool value = (_impl_._has_bits_[0] & 0x00000001u) != 0;
  PROTOBUF_ASSUME(!value || _impl_.position_ != nullptr);
  return value;
}
inline const ::pb_msg::AiVector3D& AiCamera::_internal_position() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  const ::pb_msg::AiVector3D* p = _impl_.position_;
  return p != nullptr ? *p : reinterpret_cast<const ::pb_msg::AiVector3D&>(::pb_msg::_AiVector3D_default_instance_);
}
inline const ::pb_msg::AiVector3D& AiCamera::position() const {
  // @@protoc_insertion_point(field_get:pb_msg.AiCamera.Position)
  return _internal_position();
}
inline void AiCamera::unsafe_arena_set_allocated_position(::pb_msg::AiVector3D* value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  if (GetArenaForAllocation() == nullptr) {
    delete reinterpret_cast<::google::protobuf::MessageLite*>(_impl_.position_);
  }
  _impl_.position_ = reinterpret_cast<::pb_msg::AiVector3D*>(value);
  if (value != nullptr) {
    _impl_._has_bits_[0] |= 0x00000001u;
  } else {
    _impl_._has_bits_[0] &= ~0x00000001u;
  }
  // @@protoc_insertion_point(field_unsafe_arena_set_allocated:pb_msg.AiCamera.Position)
}
inline ::pb_msg::AiVector3D* AiCamera::release_position() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);

  _impl_._has_bits_[0] &= ~0x00000001u;
  ::pb_msg::AiVector3D* released = _impl_.position_;
  _impl_.position_ = nullptr;
#ifdef PROTOBUF_FORCE_COPY_IN_RELEASE
  auto* old = reinterpret_cast<::google::protobuf::MessageLite*>(released);
  released = ::google::protobuf::internal::DuplicateIfNonNull(released);
  if (GetArenaForAllocation() == nullptr) {
    delete old;
  }
#else   // PROTOBUF_FORCE_COPY_IN_RELEASE
  if (GetArenaForAllocation() != nullptr) {
    released = ::google::protobuf::internal::DuplicateIfNonNull(released);
  }
#endif  // !PROTOBUF_FORCE_COPY_IN_RELEASE
  return released;
}
inline ::pb_msg::AiVector3D* AiCamera::unsafe_arena_release_position() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  // @@protoc_insertion_point(field_release:pb_msg.AiCamera.Position)

  _impl_._has_bits_[0] &= ~0x00000001u;
  ::pb_msg::AiVector3D* temp = _impl_.position_;
  _impl_.position_ = nullptr;
  return temp;
}
inline ::pb_msg::AiVector3D* AiCamera::_internal_mutable_position() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  _impl_._has_bits_[0] |= 0x00000001u;
  if (_impl_.position_ == nullptr) {
    auto* p = CreateMaybeMessage<::pb_msg::AiVector3D>(GetArenaForAllocation());
    _impl_.position_ = reinterpret_cast<::pb_msg::AiVector3D*>(p);
  }
  return _impl_.position_;
}
inline ::pb_msg::AiVector3D* AiCamera::mutable_position() {
  ::pb_msg::AiVector3D* _msg = _internal_mutable_position();
  // @@protoc_insertion_point(field_mutable:pb_msg.AiCamera.Position)
  return _msg;
}
inline void AiCamera::set_allocated_position(::pb_msg::AiVector3D* value) {
  ::google::protobuf::Arena* message_arena = GetArenaForAllocation();
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  if (message_arena == nullptr) {
    delete reinterpret_cast<::google::protobuf::MessageLite*>(_impl_.position_);
  }

  if (value != nullptr) {
    ::google::protobuf::Arena* submessage_arena =
        ::google::protobuf::Arena::InternalGetOwningArena(reinterpret_cast<::google::protobuf::MessageLite*>(value));
    if (message_arena != submessage_arena) {
      value = ::google::protobuf::internal::GetOwnedMessage(message_arena, value, submessage_arena);
    }
    _impl_._has_bits_[0] |= 0x00000001u;
  } else {
    _impl_._has_bits_[0] &= ~0x00000001u;
  }

  _impl_.position_ = reinterpret_cast<::pb_msg::AiVector3D*>(value);
  // @@protoc_insertion_point(field_set_allocated:pb_msg.AiCamera.Position)
}

// .pb_msg.AiVector3D Up = 3;
inline bool AiCamera::has_up() const {
  bool value = (_impl_._has_bits_[0] & 0x00000002u) != 0;
  PROTOBUF_ASSUME(!value || _impl_.up_ != nullptr);
  return value;
}
inline const ::pb_msg::AiVector3D& AiCamera::_internal_up() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  const ::pb_msg::AiVector3D* p = _impl_.up_;
  return p != nullptr ? *p : reinterpret_cast<const ::pb_msg::AiVector3D&>(::pb_msg::_AiVector3D_default_instance_);
}
inline const ::pb_msg::AiVector3D& AiCamera::up() const {
  // @@protoc_insertion_point(field_get:pb_msg.AiCamera.Up)
  return _internal_up();
}
inline void AiCamera::unsafe_arena_set_allocated_up(::pb_msg::AiVector3D* value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  if (GetArenaForAllocation() == nullptr) {
    delete reinterpret_cast<::google::protobuf::MessageLite*>(_impl_.up_);
  }
  _impl_.up_ = reinterpret_cast<::pb_msg::AiVector3D*>(value);
  if (value != nullptr) {
    _impl_._has_bits_[0] |= 0x00000002u;
  } else {
    _impl_._has_bits_[0] &= ~0x00000002u;
  }
  // @@protoc_insertion_point(field_unsafe_arena_set_allocated:pb_msg.AiCamera.Up)
}
inline ::pb_msg::AiVector3D* AiCamera::release_up() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);

  _impl_._has_bits_[0] &= ~0x00000002u;
  ::pb_msg::AiVector3D* released = _impl_.up_;
  _impl_.up_ = nullptr;
#ifdef PROTOBUF_FORCE_COPY_IN_RELEASE
  auto* old = reinterpret_cast<::google::protobuf::MessageLite*>(released);
  released = ::google::protobuf::internal::DuplicateIfNonNull(released);
  if (GetArenaForAllocation() == nullptr) {
    delete old;
  }
#else   // PROTOBUF_FORCE_COPY_IN_RELEASE
  if (GetArenaForAllocation() != nullptr) {
    released = ::google::protobuf::internal::DuplicateIfNonNull(released);
  }
#endif  // !PROTOBUF_FORCE_COPY_IN_RELEASE
  return released;
}
inline ::pb_msg::AiVector3D* AiCamera::unsafe_arena_release_up() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  // @@protoc_insertion_point(field_release:pb_msg.AiCamera.Up)

  _impl_._has_bits_[0] &= ~0x00000002u;
  ::pb_msg::AiVector3D* temp = _impl_.up_;
  _impl_.up_ = nullptr;
  return temp;
}
inline ::pb_msg::AiVector3D* AiCamera::_internal_mutable_up() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  _impl_._has_bits_[0] |= 0x00000002u;
  if (_impl_.up_ == nullptr) {
    auto* p = CreateMaybeMessage<::pb_msg::AiVector3D>(GetArenaForAllocation());
    _impl_.up_ = reinterpret_cast<::pb_msg::AiVector3D*>(p);
  }
  return _impl_.up_;
}
inline ::pb_msg::AiVector3D* AiCamera::mutable_up() {
  ::pb_msg::AiVector3D* _msg = _internal_mutable_up();
  // @@protoc_insertion_point(field_mutable:pb_msg.AiCamera.Up)
  return _msg;
}
inline void AiCamera::set_allocated_up(::pb_msg::AiVector3D* value) {
  ::google::protobuf::Arena* message_arena = GetArenaForAllocation();
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  if (message_arena == nullptr) {
    delete reinterpret_cast<::google::protobuf::MessageLite*>(_impl_.up_);
  }

  if (value != nullptr) {
    ::google::protobuf::Arena* submessage_arena =
        ::google::protobuf::Arena::InternalGetOwningArena(reinterpret_cast<::google::protobuf::MessageLite*>(value));
    if (message_arena != submessage_arena) {
      value = ::google::protobuf::internal::GetOwnedMessage(message_arena, value, submessage_arena);
    }
    _impl_._has_bits_[0] |= 0x00000002u;
  } else {
    _impl_._has_bits_[0] &= ~0x00000002u;
  }

  _impl_.up_ = reinterpret_cast<::pb_msg::AiVector3D*>(value);
  // @@protoc_insertion_point(field_set_allocated:pb_msg.AiCamera.Up)
}

// .pb_msg.AiVector3D LookAt = 4;
inline bool AiCamera::has_lookat() const {
  bool value = (_impl_._has_bits_[0] & 0x00000004u) != 0;
  PROTOBUF_ASSUME(!value || _impl_.lookat_ != nullptr);
  return value;
}
inline const ::pb_msg::AiVector3D& AiCamera::_internal_lookat() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  const ::pb_msg::AiVector3D* p = _impl_.lookat_;
  return p != nullptr ? *p : reinterpret_cast<const ::pb_msg::AiVector3D&>(::pb_msg::_AiVector3D_default_instance_);
}
inline const ::pb_msg::AiVector3D& AiCamera::lookat() const {
  // @@protoc_insertion_point(field_get:pb_msg.AiCamera.LookAt)
  return _internal_lookat();
}
inline void AiCamera::unsafe_arena_set_allocated_lookat(::pb_msg::AiVector3D* value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  if (GetArenaForAllocation() == nullptr) {
    delete reinterpret_cast<::google::protobuf::MessageLite*>(_impl_.lookat_);
  }
  _impl_.lookat_ = reinterpret_cast<::pb_msg::AiVector3D*>(value);
  if (value != nullptr) {
    _impl_._has_bits_[0] |= 0x00000004u;
  } else {
    _impl_._has_bits_[0] &= ~0x00000004u;
  }
  // @@protoc_insertion_point(field_unsafe_arena_set_allocated:pb_msg.AiCamera.LookAt)
}
inline ::pb_msg::AiVector3D* AiCamera::release_lookat() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);

  _impl_._has_bits_[0] &= ~0x00000004u;
  ::pb_msg::AiVector3D* released = _impl_.lookat_;
  _impl_.lookat_ = nullptr;
#ifdef PROTOBUF_FORCE_COPY_IN_RELEASE
  auto* old = reinterpret_cast<::google::protobuf::MessageLite*>(released);
  released = ::google::protobuf::internal::DuplicateIfNonNull(released);
  if (GetArenaForAllocation() == nullptr) {
    delete old;
  }
#else   // PROTOBUF_FORCE_COPY_IN_RELEASE
  if (GetArenaForAllocation() != nullptr) {
    released = ::google::protobuf::internal::DuplicateIfNonNull(released);
  }
#endif  // !PROTOBUF_FORCE_COPY_IN_RELEASE
  return released;
}
inline ::pb_msg::AiVector3D* AiCamera::unsafe_arena_release_lookat() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  // @@protoc_insertion_point(field_release:pb_msg.AiCamera.LookAt)

  _impl_._has_bits_[0] &= ~0x00000004u;
  ::pb_msg::AiVector3D* temp = _impl_.lookat_;
  _impl_.lookat_ = nullptr;
  return temp;
}
inline ::pb_msg::AiVector3D* AiCamera::_internal_mutable_lookat() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  _impl_._has_bits_[0] |= 0x00000004u;
  if (_impl_.lookat_ == nullptr) {
    auto* p = CreateMaybeMessage<::pb_msg::AiVector3D>(GetArenaForAllocation());
    _impl_.lookat_ = reinterpret_cast<::pb_msg::AiVector3D*>(p);
  }
  return _impl_.lookat_;
}
inline ::pb_msg::AiVector3D* AiCamera::mutable_lookat() {
  ::pb_msg::AiVector3D* _msg = _internal_mutable_lookat();
  // @@protoc_insertion_point(field_mutable:pb_msg.AiCamera.LookAt)
  return _msg;
}
inline void AiCamera::set_allocated_lookat(::pb_msg::AiVector3D* value) {
  ::google::protobuf::Arena* message_arena = GetArenaForAllocation();
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  if (message_arena == nullptr) {
    delete reinterpret_cast<::google::protobuf::MessageLite*>(_impl_.lookat_);
  }

  if (value != nullptr) {
    ::google::protobuf::Arena* submessage_arena =
        ::google::protobuf::Arena::InternalGetOwningArena(reinterpret_cast<::google::protobuf::MessageLite*>(value));
    if (message_arena != submessage_arena) {
      value = ::google::protobuf::internal::GetOwnedMessage(message_arena, value, submessage_arena);
    }
    _impl_._has_bits_[0] |= 0x00000004u;
  } else {
    _impl_._has_bits_[0] &= ~0x00000004u;
  }

  _impl_.lookat_ = reinterpret_cast<::pb_msg::AiVector3D*>(value);
  // @@protoc_insertion_point(field_set_allocated:pb_msg.AiCamera.LookAt)
}

// float HorizontalFOV = 5;
inline void AiCamera::clear_horizontalfov() {
  _impl_.horizontalfov_ = 0;
}
inline float AiCamera::horizontalfov() const {
  // @@protoc_insertion_point(field_get:pb_msg.AiCamera.HorizontalFOV)
  return _internal_horizontalfov();
}
inline void AiCamera::set_horizontalfov(float value) {
  _internal_set_horizontalfov(value);
  // @@protoc_insertion_point(field_set:pb_msg.AiCamera.HorizontalFOV)
}
inline float AiCamera::_internal_horizontalfov() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.horizontalfov_;
}
inline void AiCamera::_internal_set_horizontalfov(float value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.horizontalfov_ = value;
}

// float ClipPlaneNear = 6;
inline void AiCamera::clear_clipplanenear() {
  _impl_.clipplanenear_ = 0;
}
inline float AiCamera::clipplanenear() const {
  // @@protoc_insertion_point(field_get:pb_msg.AiCamera.ClipPlaneNear)
  return _internal_clipplanenear();
}
inline void AiCamera::set_clipplanenear(float value) {
  _internal_set_clipplanenear(value);
  // @@protoc_insertion_point(field_set:pb_msg.AiCamera.ClipPlaneNear)
}
inline float AiCamera::_internal_clipplanenear() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.clipplanenear_;
}
inline void AiCamera::_internal_set_clipplanenear(float value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.clipplanenear_ = value;
}

// float ClipPlaneFar = 7;
inline void AiCamera::clear_clipplanefar() {
  _impl_.clipplanefar_ = 0;
}
inline float AiCamera::clipplanefar() const {
  // @@protoc_insertion_point(field_get:pb_msg.AiCamera.ClipPlaneFar)
  return _internal_clipplanefar();
}
inline void AiCamera::set_clipplanefar(float value) {
  _internal_set_clipplanefar(value);
  // @@protoc_insertion_point(field_set:pb_msg.AiCamera.ClipPlaneFar)
}
inline float AiCamera::_internal_clipplanefar() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.clipplanefar_;
}
inline void AiCamera::_internal_set_clipplanefar(float value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.clipplanefar_ = value;
}

// float Aspect = 8;
inline void AiCamera::clear_aspect() {
  _impl_.aspect_ = 0;
}
inline float AiCamera::aspect() const {
  // @@protoc_insertion_point(field_get:pb_msg.AiCamera.Aspect)
  return _internal_aspect();
}
inline void AiCamera::set_aspect(float value) {
  _internal_set_aspect(value);
  // @@protoc_insertion_point(field_set:pb_msg.AiCamera.Aspect)
}
inline float AiCamera::_internal_aspect() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.aspect_;
}
inline void AiCamera::_internal_set_aspect(float value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.aspect_ = value;
}

// float OrthographicWidth = 9;
inline void AiCamera::clear_orthographicwidth() {
  _impl_.orthographicwidth_ = 0;
}
inline float AiCamera::orthographicwidth() const {
  // @@protoc_insertion_point(field_get:pb_msg.AiCamera.OrthographicWidth)
  return _internal_orthographicwidth();
}
inline void AiCamera::set_orthographicwidth(float value) {
  _internal_set_orthographicwidth(value);
  // @@protoc_insertion_point(field_set:pb_msg.AiCamera.OrthographicWidth)
}
inline float AiCamera::_internal_orthographicwidth() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.orthographicwidth_;
}
inline void AiCamera::_internal_set_orthographicwidth(float value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.orthographicwidth_ = value;
}

#ifdef __GNUC__
#pragma GCC diagnostic pop
#endif  // __GNUC__

// @@protoc_insertion_point(namespace_scope)
}  // namespace pb_msg


// @@protoc_insertion_point(global_scope)

#include "google/protobuf/port_undef.inc"

#endif  // GOOGLE_PROTOBUF_INCLUDED_common_2fpb_5fmsg_2fai_5fcamera_2eproto_2epb_2eh
