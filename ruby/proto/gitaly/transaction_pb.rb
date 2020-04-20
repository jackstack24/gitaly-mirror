# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: transaction.proto

require 'google/protobuf'

require 'lint_pb'
require 'shared_pb'
Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("transaction.proto", :syntax => :proto3) do
    add_message "gitaly.RegisterTransactionRequest" do
      optional :repository, :message, 1, "gitaly.Repository"
      repeated :nodes, :string, 2
    end
    add_message "gitaly.RegisterTransactionResponse" do
      optional :transaction_id, :uint64, 1
    end
    add_message "gitaly.StartTransactionRequest" do
      optional :repository, :message, 1, "gitaly.Repository"
      optional :transaction_id, :uint64, 2
      optional :node, :string, 3
      optional :reference_updates_hash, :bytes, 4
    end
    add_message "gitaly.StartTransactionResponse" do
      optional :state, :enum, 1, "gitaly.StartTransactionResponse.TransactionState"
    end
    add_enum "gitaly.StartTransactionResponse.TransactionState" do
      value :COMMIT, 0
      value :ABORT, 1
    end
    add_message "gitaly.CancelTransactionRequest" do
      optional :repository, :message, 1, "gitaly.Repository"
      optional :transaction_id, :uint64, 2
    end
    add_message "gitaly.CancelTransactionResponse" do
    end
  end
end

module Gitaly
  RegisterTransactionRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitaly.RegisterTransactionRequest").msgclass
  RegisterTransactionResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitaly.RegisterTransactionResponse").msgclass
  StartTransactionRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitaly.StartTransactionRequest").msgclass
  StartTransactionResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitaly.StartTransactionResponse").msgclass
  StartTransactionResponse::TransactionState = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitaly.StartTransactionResponse.TransactionState").enummodule
  CancelTransactionRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitaly.CancelTransactionRequest").msgclass
  CancelTransactionResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitaly.CancelTransactionResponse").msgclass
end
